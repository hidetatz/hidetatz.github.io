import os
github_context = os.environ.get("GITHUB_CONTEXT", "{}")
print(github_context)
