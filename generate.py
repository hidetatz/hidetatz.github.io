import dataclasses
import datetime
import os
import os.path
import shutil
import string
import xml.etree.ElementTree as ET

import markdown

import template

md = markdown.Markdown(extensions=["tables", "fenced_code"])

def new_article(filename):
    os.makedirs("data/articles", exist_ok=True)
    with open(f"data/articles/{filename}.md", "w") as f:
        f.write(f"""title: {filename}
timestamp: {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}
url: necessary for external link
lang: ja/en
---
""")
    return f"data/articles/{filename}.md"

def new_diary(title, ts, body):
    os.makedirs("data/diaries", exist_ok=True)
    filename = f"data/diaries/{ts.strftime('%Y-%m-%d')}.md"
    with open(filename, "w") as f:
        f.write(f"""title: {title}
timestamp: {ts.strftime('%Y-%m-%d %H:%M:%S')}
lang: ja
---
{body}
---
""")
    return filename

class Sitemap:
    def __init__(self, articles, diaries):
        self.articles = articles
        self.diaries = diaries

        if len(diaries) == 0:
            self.latest_ts = articles[0].ts_atom()
        elif articles[0].timestamp > diaries[0].timestamp:
            self.latest_ts = articles[0].ts_atom()
        else:
            self.latest_ts = diaries[0].ts_atom()

    def save_xml(self, out):
        urlset = ET.Element("urlset")
        urlset.set("xmlns", "http://www.sitemaps.org/schemas/sitemap/0.9")

        url_elem = ET.SubElement(urlset, "url")
        loc = ET.SubElement(url_elem, "loc")
        loc.text = f"https://hidetatz.github.io"
        lastmod = ET.SubElement(url_elem, "lastmod")
        lastmod.text = self.latest_ts

        for article in self.articles:
            url_elem = ET.SubElement(urlset, "url")
            loc = ET.SubElement(url_elem, "loc")
            if article.external_url == "":
                loc.text = f"https://hidetatz.github.io/{article.url_path('articles')}"
            else:
                loc.text = article.external_url
            lastmod = ET.SubElement(url_elem, "lastmod")
            lastmod.text = article.ts_atom()

        for diary in self.diaries:
            url_elem = ET.SubElement(urlset, "url")
            loc = ET.SubElement(url_elem, "loc")
            loc.text = f"https://hidetatz.github.io/{diary.url_path('diary')}"
            lastmod = ET.SubElement(url_elem, "lastmod")
            lastmod.text = diary.ts_atom()
        
        tree = ET.ElementTree(element=urlset)
        tree.write(out, encoding="utf-8", xml_declaration=True)

class AtomFeed:
    def __init__(self, articles, diaries):
        articles_tup = [("article", article) for article in articles]
        diaries_tup = [("diary", diary) for diary in diaries]

        self.all_articles = articles_tup + diaries_tup
        self.all_articles.sort(key=lambda a: a[1].timestamp, reverse=True)

    def save_xml(self, out):
        feed = ET.Element('feed')
        feed.set("xmlns", "http://www.w3.org/2005/Atom")
        title = ET.SubElement(feed, "title")
        title.text = "hidetatz.github.io | Hidetatz Yaginuma"

        _id = ET.SubElement(feed, "id")
        _id.text = "https://hidetatz.github.io"

        updated = ET.SubElement(feed, "updated")
        updated.text = self.all_articles[0][1].ts_atom()

        link = ET.SubElement(feed, "link")
        link.set("href", "https://hidetatz.github.io")

        author = ET.SubElement(feed, "author")
        author_name = ET.SubElement(author, "name")
        author_name.text = "Hidetatz Yaginuma"
        author_email = ET.SubElement(author, "email")
        author_email.text = "hidetatz@gmail.com"
        
        for i in range(20):
            ent = self.all_articles[i]
            typ, article = ent[0], ent[1]

            article_url = article.external_url
            if article.external_url == "":
                article_url = "https://hidetatz.github.io/" + article.url_path("articles" if typ == "article" else "diary")

            entry = ET.SubElement(feed, "entry")

            entry_title = ET.SubElement(entry, "title")
            entry_title.text = article.title

            entry_updated = ET.SubElement(entry, "updated")
            entry_updated.text = article.ts_atom()

            entry_id = ET.SubElement(entry, "id")
            entry_id.text = article_url

            entry_link = ET.SubElement(entry, "link")
            entry_link.set("href", article_url)
            entry_link.set("rel", "alternate")

            entry_summary = ET.SubElement(entry, "summary")
            entry_summary.set("type", "html")
            entry_summary.text = "The post first appeared on hidetatz.github.io."

            entry_author = ET.SubElement(entry, "author")
            entry_author_name = ET.SubElement(entry_author, "name")
            entry_author_name.text = "Hidetatz Yaginuma"
            entry_author_email = ET.SubElement(entry_author, "email")
            entry_author_email.text = "hidetatz@gmail.com"

        tree = ET.ElementTree(element=feed)
        tree.write(out, encoding='utf-8', xml_declaration=True)

@dataclasses.dataclass
class MarkdownArticle:
    filename_no_ext: str
    title: str
    timestamp: datetime.time
    content: list
    external_url: str = ""
    lang: str = "en"

    def __init__(self, file):
        self.filename_no_ext, _ = os.path.splitext(os.path.basename(file.name))
        lines = file.read().splitlines()
        self.content = []
        in_front_matter = True
        for line in lines:
            if in_front_matter:
                # assume every article must have yaml front matter finishing line
                if line == "---":
                    in_front_matter = False
                    continue

                key, val = line.split(": ")

                if key == "timestamp":
                    self.timestamp = datetime.datetime.strptime(val, "%Y-%m-%d %H:%M:%S")
                    continue

                if key == "url":
                    self.external_url = val
                    continue

                if key == "lang":
                    self.lang = val
                    continue

                if key == "title":
                    self.title = val
                    continue

            else:
                self.content.append(line)

    def ts_display(self): 
        return self.timestamp.strftime('%Y/%m/%d')

    def ts_atom(self): 
        return self.timestamp.strftime('%Y-%m-%dT%H:%M:%SZ')

    def url_path(self, directory): 
        return f"{directory}/{self.ts_display()}/{self.filename_no_ext}"

    def md_link(self, directory):
        if self.external_url != "":
            return f"[{self.title}]({self.external_url})"

        url_path = self.url_path(directory)
        return f"/[{self.title}]({url_path})"

    def as_html(self, x=True):
        if self.external_url != "":
            return

        t = string.Template(template.article_content)
        content = t.substitute(title=self.title, timestamp=self.ts_display(), content='\n'.join(self.content))
        content = md.convert(content)
        if x:
            content += '\n<p><a href="https://twitter.com/share?ref_src=twsrc%5Etfw" class="twitter-share-button" data-via="hidetatz" data-show-count="false">Tweet</a><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></p>'
        return string.Template(template.html_page).substitute(title=self.title, body=content)

@dataclasses.dataclass
class BlogSite:
    articles: list
    diaries: list

    def save(self, location, content):
        os.makedirs(os.path.dirname(location), exist_ok=True)
        with open(location, "w") as f:
            f.write(content)

    def copy(self, filename):
        src = f"data/{filename}"
        dst = f"docs/{filename}"
        shutil.copyfile(src, dst)

    def tmpl_md_as_html(self, title, tmpl, **kwargs):
        tmpl = string.Template(tmpl)
        body = tmpl.substitute(**kwargs)
        body = md.convert(body)
        return string.Template(template.html_page).substitute(title=title, body=body)

    def generate(self):
        shutil.rmtree("docs")
        os.mkdir("docs")

        # Copy static files from data to docs
        self.copy("robots.txt")
        self.copy("markdown.css")
        self.copy("syntax.css")
        self.copy("highlight.pack.js")
        self.copy("favicon.ico")

        # Create articles and article index. Article index is also the whole site index.
        article_links = []
        article_ja_links = []

        for article in self.articles:
            target = article_ja_links if article.lang == "ja" else article_links
            target.append(f"{article.ts_display()} - {article.md_link('articles')}  ")

            if article.external_url == "":
                self.save(f"docs/{article.url_path('articles')}/index.html", article.as_html())

        index = self.tmpl_md_as_html("hidetatz.github.io", template.index_page_md, articles="\n".join(article_links), articles_ja="\n".join(article_ja_links))
        self.save(f"docs/index.html", index)

        # Create diaries and diaries index.
        diary_links = []

        for diary in self.diaries:
            diary_links.append(f"{diary.md_link('diary')}  ")
            self.save(f"docs/{diary.url_path('diary')}/index.html", diary.as_html(False))

        diary_index = self.tmpl_md_as_html("diary | hidetatz.github.io", template.diary_index_page_md, diaries="\n".join(diary_links))
        self.save(f"docs/diary/index.html", diary_index)

        # Generate 404 page.
        not_found = self.tmpl_md_as_html("404 | hidetatz.github.io", template.not_found_page_md, recent_articles="\n".join(article_links[:5]))
        self.save(f"docs/404.html", not_found)

        # Generate sitemap and rss feed.
        sitemap = Sitemap(self.articles, self.diaries)
        sitemap.save_xml("docs/sitemap.xml")

        atomfeed = AtomFeed(self.articles, self.diaries)
        atomfeed.save_xml("docs/feed.xml")

def read_article_files(directory):
    files = [f for f in os.listdir(directory) if os.path.isfile(f"{directory}/{f}")]
    articles = []
    for file in files:
        with open(f"{directory}/{file}") as f:
            articles.append(MarkdownArticle(f))

    articles.sort(key=lambda article: article.timestamp, reverse=True)

    return articles

def generate():
    articles = read_article_files("data/articles")
    diaries = read_article_files("data/diaries")

    BlogSite(articles, diaries).generate()

if __name__ == "__main__":
    generate()
