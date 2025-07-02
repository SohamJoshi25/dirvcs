package data

import "path"

var BASE_PATH string = "./.dirvcs"
var CONFIG_PATH string = path.Join(BASE_PATH, "config.yaml")
var LOGS_PATH string = path.Join(BASE_PATH, "logs.json")
var IGNORE_PATH string = path.Join(BASE_PATH, ".ignore")
var TREES_PATH string = path.Join(BASE_PATH, "trees")
var TREE_LOG_PATH string = path.Join(TREES_PATH, "treelogs.json")
