import os

import mistune
import requests
import yaml

# import pprint


# pp = pprint.PrettyPrinter(indent=1)


class IndentDumper(yaml.SafeDumper):
    def ignore_aliases(self, _):
        return True

    def increase_indent(self, flow=False, indentless=False):
        return super(IndentDumper, self).increase_indent(flow, False)


def parse_markdown(md):
    heading = md[0]
    command = heading["children"][0]["text"]
    r = {}

    for obj in md:
        if obj["type"] == "list":
            r["name"] = command
            text_children = obj["children"][0]["children"][0]["children"]
            text = ""
            for t in text_children:
                text += t["text"]
            r["description"] = text.rstrip(":")
        elif obj["type"] == "paragraph":
            r["example"] = obj["children"][0]["text"]
            if "name" in r:
                yield r
                r = {}
            else:
                print("error: found paragraph type without matching list")


def filename(outdir, section):
    name = section.replace(" ", "-").lower()
    return f"{outdir}/{name}.yml"


def write_yaml(outdir, data):
    for section, commands in data.items():
        out = {
            "section": section,
            "commands": commands,
        }
        with open(filename(outdir, section), "w") as f:
            yaml.dump(out, f, Dumper=IndentDumper, sort_keys=False)


class Parser:
    def __init__(self, path, pages, outdir):
        self.tldr_url = f"https://raw.githubusercontent.com/tldr-pages/tldr/main/{pages}"
        self.path = path
        self.outdir = outdir
        self.markdown = mistune.create_markdown(renderer="ast")

    def get_page(self, page_url):
        res = requests.get(f"{self.tldr_url}/{page_url}")
        if res.status_code != 200:
            print(f"error: failed to get page: {page_url}")
        return res.text

    def run(self):
        yaml_data = {}

        with open(self.path) as f:
            config = yaml.safe_load(f)
            for entry in config:
                section = entry["section"]
                if section not in yaml_data:
                    yaml_data[section] = []
                for page in entry["pages"]:
                    os, _ = page["path"].split("/")
                    filters = page["filters"] if "filters" in page else []
                    if os != "common":
                        filters.append(os)
                    text = self.get_page(page["path"])
                    md = self.markdown(text)
                    for command in parse_markdown(md):
                        if len(filters) > 0:
                            command["filters"] = filters
                        yaml_data[section].append(command)

        write_yaml(self.outdir, yaml_data)


p = os.path.dirname(__file__)
Parser(f"{p}/tldr_pages.yml", "pages", f"{p}/../data").run()
