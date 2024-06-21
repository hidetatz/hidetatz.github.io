import datetime
import os
import sys

def new_article(filename):
    os.makedirs("articles", exist_ok=True)
    with open(f"articles/{filename}.md", "w") as f:
        f.write(f"""title: {filename}
timestamp: {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}
lang: ja/en
---
""")
    return f"articles/{filename}.md"

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("filename must be given!")

    else:
        filename = sys.argv[1]
        print(new_article(filename))
