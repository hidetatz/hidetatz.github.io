import datetime
import os
import re

from github import Github, Auth

import generate

class dotdict(dict):
    __getattr__ = dict.get
    __setattr__ = dict.__setitem__
    __delattr__ = dict.__delitem__

ctx = dotdict(json.loads(os.environ.get("GITHUB_CONTEXT")))
gh = Github(auth=Auth.Token(os.environ.get("GITHUB_TOKEN")))

# closed: do publish
if ctx.event.action == "closed":
    repo = gh.get_repo("hidetatz/hidetatz.github.io")
    issue = repo.get_issue(ctx.event.issue.number)

    ts = issue.closed_at
    pattern = "20\d{2}年\d{1,2}月\d{1,2}日"
    result = re.search(pattern, issue.title)
    if result:
        ts = datetime.datetime.strptime(result.group(), "%Y年%m月%d日")
    generate.new_diary(issue.title, ts, issue.body)
    generate.generate()
    gh.close()

# # edited: if the issue is already closed, do publish
# if ctx["event"]["action"] == "edited":
#     # publish
#     gh.close()
#     return
# 
# # reopened, deleted: unpublish
# if ctx["event"]["action"] in ["reopened", "deleted"]:
#     # unpublish
#     gh.close()
#     return

else:
    print(f"event {ctx['event']['action']} is ok not to handle")
    gh.close()
