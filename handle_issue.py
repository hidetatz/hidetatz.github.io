import datetime
import json
import os
import re
import subprocess

from github import Github, Auth

import generate

ctx = json.loads(os.environ.get("GITHUB_CONTEXT"))
gh_token = os.environ.get("GITHUB_TOKEN")
gh = Github(auth=Auth.Token(gh_token))

if ctx["event"]["action"] == "closed":
    repo = gh.get_repo("hidetatz/hidetatz.github.io")
    issue = repo.get_issue(ctx["event"]["issue"]["number"])

    ts = issue.closed_at
    pattern = "20\d{2}年\d{1,2}月\d{1,2}日"
    result = re.search(pattern, issue.title)
    if result:
        ts = datetime.datetime.strptime(result.group(), "%Y年%m月%d日")
    generate.new_diary(issue.title, ts, issue.body)
    generate.generate()
    gh.close()

    subprocess.run(["git", "remote", "add", "origin" f"https://hidetatz:{gh_token}@github.com/hidetatz/hidetatz.github.io.git"])
    subprocess.run(["git", "config", "--global", "user.email", "hidetatz@gmail.com"])
    subprocess.run(["git", "config", "--global", "user.name", "Hidetatz Yaginuma in CI"])
    subprocess.run(["git", "add", "docs"])
    subprocess.run(["git", "add", "data"])
    subprocess.run(["git", "commit", "-m", "update diary"])
    subprocess.run(["git", "push", "origin", "master"])

else:
    print(f"event {ctx['event']['action']} is ok not to handle")
    gh.close()
