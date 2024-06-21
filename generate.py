import dataclasses
import datetime
import http.client
import json
import os
import os.path
import re
import shutil
import string
import subprocess
import xml.etree.ElementTree as ET

import markdown

import template

md = markdown.Markdown(extensions=["tables", "fenced_code"])

def new_article(filename):
    os.makedirs("data/articles", exist_ok=True)
    with open(f"data/articles/{filename}.md", "w") as f:
        f.write(f"""title: {filename}
timestamp: {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}
lang: ja/en
---
""")
    return f"data/articles/{filename}.md"

class Sitemap:
    def __init__(self, articles, diaries):
        self.articles = articles
        self.diaries = diaries

        # self.latest_ts = diaries[0].ts_iso8601()
        self.latest_ts = articles[0].ts_iso8601()

    def save_xml(self, out):
        urlset = ET.Element("urlset")
        urlset.set("xmlns", "http://www.sitemaps.org/schemas/sitemap/0.9")

        url_elem = ET.SubElement(urlset, "url")
        loc = ET.SubElement(url_elem, "loc")
        loc.text = f"https://hidetatz.github.io"
        lastmod = ET.SubElement(url_elem, "lastmod")
        lastmod.text = self.latest_ts

        for entry in self.articles + self.diaries:
            url_elem = ET.SubElement(urlset, "url")
            loc = ET.SubElement(url_elem, "loc")
            loc.text = f"https://hidetatz.github.io/{entry.url_path()}"
            lastmod = ET.SubElement(url_elem, "lastmod")
            lastmod.text = entry.ts_iso8601()
        
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
        updated.text = self.all_articles[0][1].ts_iso8601()

        link = ET.SubElement(feed, "link")
        link.set("href", "https://hidetatz.github.io")

        author = ET.SubElement(feed, "author")
        author_name = ET.SubElement(author, "name")
        author_name.text = "Hidetatz Yaginuma"
        author_email = ET.SubElement(author, "email")
        author_email.text = "hidetatz@gmail.com"
        
        for i in range(15):
            ent = self.all_articles[i]
            typ, article = ent[0], ent[1]

            article_url = f"https://hidetatz.github.io/{article.url_path()}"

            entry = ET.SubElement(feed, "entry")

            entry_title = ET.SubElement(entry, "title")
            entry_title.text = article.title

            entry_updated = ET.SubElement(entry, "updated")
            entry_updated.text = article.ts_iso8601()

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

class Entry:
    def __init__(self, title, timestamp, content):
        self.title = title
        self.timestamp = timestamp
        self.content = content

    def ts_display(self): 
        return self.timestamp.strftime('%Y/%m/%d')

    def ts_short(self): 
        return self.timestamp.strftime('%Y%m%d')

    def ts_iso8601(self): 
        return self.timestamp.strftime('%Y-%m-%dT%H:%M:%SZ')

class Article(Entry):
    def __init__(self, filepath):
        with open(filepath) as f:
            self.filename_no_ext, _ = os.path.splitext(os.path.basename(f.name))
            lines = f.read().splitlines()

        self.content = []
        in_front_matter = True
        for line in lines:
            if in_front_matter:
                if line == "---":
                    in_front_matter = False
                    continue

                key, val = line.split(": ")

                if key == "timestamp":
                    timestamp = datetime.datetime.strptime(val, "%Y-%m-%d %H:%M:%S")

                elif key == "title":
                    title = val

                elif key == "lang":
                    self.lang = val

            else:
                self.content.append(line)

        super().__init__(title, timestamp, "\n".join(self.content))

    def url_path(self): 
        return f"/{self.filename_no_ext}"

    def to_html(self):
        t = string.Template(template.article_content)
        content = md.convert(t.substitute(title=self.title, timestamp=self.ts_display(), content=self.content))
        content += '\n<p><a href="https://twitter.com/share?ref_src=twsrc%5Etfw" class="twitter-share-button" data-via="hidetatz" data-show-count="false">Tweet</a><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></p>'
        return content

class Diary(Entry):
    def __init__(self, issue):
        ts = datetime.datetime.strptime(issue["created_at"], "%Y-%m-%dT%H:%M:%SZ")
        pattern = "20\d{2}年\d{1,2}月\d{1,2}日"
        result = re.search(pattern, issue["title"])
        if result:
            ts = datetime.datetime.strptime(result.group(), "%Y年%m月%d日")

        super().__init__(issue["title"], ts, issue["body"])

    def url_path(self): 
        return f"/{self.ts_short()}"

    def to_html(self):
        t = string.Template(template.diary_content)
        return md.convert(t.substitute(title=self.title, content=self.content))

class Blog:
    def __init__(self, root, gh_token):
        self.root = root
        self.gh_token = gh_token

    def save(self, location, content):
        os.makedirs(os.path.dirname(location), exist_ok=True)
        with open(location, "w") as f:
            f.write(content)

    def copy(self, filename):
        src = f"static/{filename}"
        dst = f"{self.root}/{filename}"
        os.makedirs(os.path.dirname(dst), exist_ok=True)
        shutil.copyfile(src, dst)

    def to_html(self, title, body):
        return string.Template(template.html_page).substitute(title=title, body=body)

    def tmpl_md_as_html(self, title, tmpl, **kwargs):
        return self.to_html(title, md.convert(string.Template(tmpl).substitute(**kwargs)))

    def create_articles(self):
        articles = []
        for file in os.listdir("articles"):
            article = Article(f"articles/{file}")
            self.save(f"{self.root}/{article.url_path()}/index.html", self.to_html(article.title, article.to_html()))
            articles.append(article)

        articles.sort(key=lambda article: article.timestamp, reverse=True)
        return articles

    def create_diaries(self):
        def fetch_all_issues(gh_token):
            issues = []
            page=1
            while True:
                conn = http.client.HTTPSConnection("api.github.com")
                headers = {
                    "Accept": "application/vnd.github+json",
                    "Authorization": f"Bearer {gh_token}",
                    "X-GitHub-Api-Version": "2022-11-28",
                    "User-Agent": "hidetatz.github.io",
                }
                conn.request("GET", "/repos/hidetatz/hidetatz.github.io/issues?state=open&creator=hidetatz&per_page=100&page={page}", headers=headers)
                resp = conn.getresponse()
                body = json.loads(resp.read().decode("utf-8"))
                issues += body

                if "Link" not in resp.headers:
                    break
                
                if 'relname="next"' not in resp.headers["Link"]:
                    break

                page += 1

            return issues

        
        diaries = []
        for issue in fetch_all_issues(self.gh_token):
            diary = Diary(issue)
            self.save(f"{self.root}/{diary.url_path()}/index.html", self.to_html(diary.title, diary.to_html()))
            diaries.append(diary)

        diaries.sort(key=lambda diary: diary.timestamp, reverse=True)
        return diaries 

    def generate_gh_pages(self):
        shutil.rmtree(self.root, ignore_errors=True)

        self.copy("robots.txt")
        self.copy("markdown.css")
        self.copy("syntax.css")
        self.copy("highlight.pack.js")
        self.copy("favicon.ico")

        # Create each article/diary pages.

        articles = self.create_articles()
        diaries = self.create_diaries()

        # Create index pages.

        article_links = []
        for article in articles:
            article_links.append(f"{article.ts_display()} - [{article.title}]({article.url_path()})  ")

        index = self.tmpl_md_as_html("hidetatz.github.io", template.index_page_md, articles="\n".join(article_links))
        self.save(f"{self.root}/index.html", index)

        diary_links = []
        for diary in diaries:
            diary_links.append(f"[{diary.title}]({diary.url_path()})  ")

        diary_index = self.tmpl_md_as_html("diary | hidetatz.github.io", template.diary_index_page_md, diaries="\n".join(diary_links))
        self.save(f"{self.root}/diary/index.html", diary_index)

        # Generate 404 page.
        not_found = self.tmpl_md_as_html("404 | hidetatz.github.io", template.not_found_page_md, recent_articles="\n".join(article_links[:5]))
        self.save(f"{self.root}/404.html", not_found)

        # Generate sitemap and rss feed.
        Sitemap(articles, diaries).save_xml(f"{self.root}/sitemap.xml")
        AtomFeed(articles, diaries).save_xml(f"{self.root}/feed.xml")

    def generate_and_push(self):
        self.generate_gh_pages()

        shutil.move(f"./{self.root}", f"../{self.root}")

        # subprocess.run(["git", "remote", "add", "origin" f"https://hidetatz:{self.gh_token}@github.com/hidetatz/hidetatz.github.io.git"])
        subprocess.run(["git", "config", "--global", "user.email", "hidetatz@gmail.com"])
        subprocess.run(["git", "config", "--global", "user.name", "Hidetatz Yaginuma in CI"])

        subprocess.run(["git", "checkout", "-b", "gh-pages"])

        for item in os.listdir("."):
            if item == '.git':
                continue
            if os.path.isfile(item):
                os.remove(item)
            elif os.path.isdir(item):
                shutil.rmtree(item)

        for file in glob.glob(f"../{self.root}/*")
            shutil.move(file, f"./")

        subprocess.run(["git", "add", "."])
        subprocess.run(["git", "commit", "-m", "update"])
        subprocess.run(["git", "push", "origin", "gh-pages"])


if __name__ == "__main__":
    Blog("gh_pages", os.environ.get("GITHUB_TOKEN")).generate_and_push()
