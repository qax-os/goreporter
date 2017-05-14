package tools

import (
	"testing"
)

const structJson = `{
    "project": "github.com\\wgliang\\logcool",
    "score": 0,
    "grade": 0,
    "metrics": {
        "CopyCodeTips": {
            "name": "CopyCode",
            "description": "Query all duplicate code in the project and give duplicate code locations and rows.",
            "summaries": {
                "\u0000": {
                    "name": "2",
                    "description": "",
                    "errors": [
                        {
                            "line_number": 22,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\cmd\\cmd.go:104,125\n"
                        },
                        {
                            "line_number": 22,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:449,470\n"
                        }
                    ]
                },
                "\u0001": {
                    "name": "6",
                    "description": "",
                    "errors": [
                        {
                            "line_number": 16,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\filter\\grok\\grok.go:26,41\n"
                        },
                        {
                            "line_number": 17,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\filter\\split\\split.go:26,42\n"
                        },
                        {
                            "line_number": 16,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\filter\\zeus\\zeus.go:25,40\n"
                        },
                        {
                            "line_number": 15,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\output\\email\\email.go:33,47\n"
                        },
                        {
                            "line_number": 15,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\output\\lexec\\lexec.go:26,40\n"
                        },
                        {
                            "line_number": 15,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\output\\stdout\\stdout.go:25,39\n"
                        }
                    ]
                },
                "\u0002": {
                    "name": "2",
                    "description": "",
                    "errors": [
                        {
                            "line_number": 29,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\utils\\filter.go:60,88\n"
                        },
                        {
                            "line_number": 29,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\utils\\output.go:63,91\n"
                        }
                    ]
                },
                "\u0003": {
                    "name": "2",
                    "description": "",
                    "errors": [
                        {
                            "line_number": 19,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:110,128\n"
                        },
                        {
                            "line_number": 19,
                            "error_string": "..\\src\\github.com\\wgliang\\logcool\\input\\stdin\\stdin.go:32,50\n"
                        }
                    ]
                }
            },
            "weight": 0.1,
            "percentage": 90,
            "error": ""
        },
        "CycloTips": {
            "name": "Cyclo",
            "description": "Computing all [.go] file's cyclo,and as an important indicator of the quality of the code.",
            "summaries": {
                "github.com\\wgliang\\logcool": {
                    "name": "github.com\\wgliang\\logcool",
                    "description": "8.00",
                    "errors": [
                        {
                            "line_number": 8,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\logcool.go:35:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\cmd": {
                    "name": "github.com\\wgliang\\logcool\\cmd",
                    "description": "2.23",
                    "errors": [
                        {
                            "line_number": 5,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd.go:67:1"
                        },
                        {
                            "line_number": 5,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd.go:44:1"
                        },
                        {
                            "line_number": 5,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd.go:104:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd.go:33:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd.go:22:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd_test.go:39:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd.go:95:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd.go:85:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd_test.go:47:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd_test.go:35:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd_test.go:31:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd_test.go:8:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\cmd\\cmd_test.go:51:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\filter\\grok": {
                    "name": "github.com\\wgliang\\logcool\\filter\\grok",
                    "description": "2.17",
                    "errors": [
                        {
                            "line_number": 5,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\grok\\grok.go:44:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\grok\\grok_test.go:25:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\grok\\grok_test.go:16:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\grok\\grok.go:26:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\grok\\grok_test.go:12:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\grok\\grok.go:21:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\filter\\metrics": {
                    "name": "github.com\\wgliang\\logcool\\filter\\metrics",
                    "description": "2.00",
                    "errors": [
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\metrics\\metrics.go:50:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\metrics\\metrics_test.go:25:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\metrics\\metrics_test.go:16:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\metrics\\metrics.go:30:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\metrics\\metrics_test.go:12:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\metrics\\metrics.go:25:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\filter\\split": {
                    "name": "github.com\\wgliang\\logcool\\filter\\split",
                    "description": "1.83",
                    "errors": [
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\split\\split.go:45:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\split\\split_test.go:25:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\split\\split_test.go:16:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\split\\split.go:26:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\split\\split_test.go:12:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\split\\split.go:21:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\filter\\zeus": {
                    "name": "github.com\\wgliang\\logcool\\filter\\zeus",
                    "description": "1.83",
                    "errors": [
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\zeus\\zeus.go:43:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\zeus\\zeus_test.go:25:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\zeus\\zeus_test.go:16:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\zeus\\zeus.go:25:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\zeus\\zeus_test.go:12:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\filter\\zeus\\zeus.go:20:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\input\\collectd": {
                    "name": "github.com\\wgliang\\logcool\\input\\collectd",
                    "description": "1.50",
                    "errors": [
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:136:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:110:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd_test.go:24:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd_test.go:15:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:184:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:173:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:239:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:131:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:198:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd_test.go:11:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:280:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:233:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:105:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\collectd\\collectd.go:214:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\input\\file": {
                    "name": "github.com\\wgliang\\logcool\\input\\file",
                    "description": "4.89",
                    "errors": [
                        {
                            "line_number": 17,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:222:1"
                        },
                        {
                            "line_number": 11,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:403:1"
                        },
                        {
                            "line_number": 10,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:167:1"
                        },
                        {
                            "line_number": 7,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:369:1"
                        },
                        {
                            "line_number": 6,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:84:1"
                        },
                        {
                            "line_number": 5,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:449:1"
                        },
                        {
                            "line_number": 5,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:114:1"
                        },
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:141:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:56:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:328:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:158:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file_test.go:24:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:355:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:340:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file_test.go:15:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:51:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file_test.go:11:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\file\\file.go:79:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\input\\http": {
                    "name": "github.com\\wgliang\\logcool\\input\\http",
                    "description": "2.33",
                    "errors": [
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\http\\http.go:84:1"
                        },
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\http\\http.go:65:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\http\\http_test.go:13:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\http\\http.go:35:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\http\\http_test.go:32:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\http\\http_test.go:41:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\http\\http.go:60:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\http\\http.go:31:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\http\\http_test.go:28:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\input\\stdin": {
                    "name": "github.com\\wgliang\\logcool\\input\\stdin",
                    "description": "2.00",
                    "errors": [
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\stdin\\stdin.go:57:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\stdin\\stdin.go:32:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\stdin\\stdin_test.go:15:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\stdin\\stdin_test.go:24:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\stdin\\stdin.go:27:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\stdin\\stdin_test.go:11:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\input\\stdin\\stdin.go:53:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\output\\email": {
                    "name": "github.com\\wgliang\\logcool\\output\\email",
                    "description": "1.86",
                    "errors": [
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\email\\email.go:55:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\email\\email_test.go:16:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\email\\email.go:33:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\email\\email_test.go:34:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\email\\email.go:28:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\email\\email_test.go:12:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\email\\email.go:50:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\output\\lexec": {
                    "name": "github.com\\wgliang\\logcool\\output\\lexec",
                    "description": "1.83",
                    "errors": [
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\lexec\\lexec.go:43:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\lexec\\lexec_test.go:29:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\lexec\\lexec_test.go:16:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\lexec\\lexec.go:26:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\lexec\\lexec_test.go:12:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\lexec\\lexec.go:21:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\output\\redis": {
                    "name": "github.com\\wgliang\\logcool\\output\\redis",
                    "description": "2.44",
                    "errors": [
                        {
                            "line_number": 7,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\redis\\redis.go:72:1"
                        },
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\redis\\redis.go:106:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\redis\\redis.go:38:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\redis\\redis.go:65:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\redis\\redis_test.go:36:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\redis\\redis_test.go:17:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\redis\\redis.go:60:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\redis\\redis.go:33:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\redis\\redis_test.go:13:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\output\\stdout": {
                    "name": "github.com\\wgliang\\logcool\\output\\stdout",
                    "description": "1.67",
                    "errors": [
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\stdout\\stdout_test.go:29:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\stdout\\stdout_test.go:16:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\stdout\\stdout.go:42:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\stdout\\stdout.go:25:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\stdout\\stdout_test.go:12:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\output\\stdout\\stdout.go:20:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\utils": {
                    "name": "github.com\\wgliang\\logcool\\utils",
                    "description": "2.28",
                    "errors": [
                        {
                            "line_number": 8,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:160:1"
                        },
                        {
                            "line_number": 7,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\input.go:52:1"
                        },
                        {
                            "line_number": 7,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\filter.go:60:1"
                        },
                        {
                            "line_number": 7,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\output.go:63:1"
                        },
                        {
                            "line_number": 6,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\output.go:40:1"
                        },
                        {
                            "line_number": 5,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\filter.go:39:1"
                        },
                        {
                            "line_number": 5,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:62:1"
                        },
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:49:1"
                        },
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:187:1"
                        },
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:74:1"
                        },
                        {
                            "line_number": 4,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:85:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\input.go:39:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:118:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:76:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:143:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:131:1"
                        },
                        {
                            "line_number": 3,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:102:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:209:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent_test.go:42:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\input_test.go:14:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:29:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config_test.go:67:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:118:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\filter_test.go:14:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent_test.go:26:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\output_test.go:14:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:98:1"
                        },
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:88:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent_test.go:12:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config_test.go:63:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config_test.go:19:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config_test.go:84:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config_test.go:50:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:37:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config_test.go:33:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config_test.go:37:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:65:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\filter.go:28:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent_test.go:95:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent_test.go:56:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\filter_test.go:10:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:70:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\output.go:29:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\filter.go:33:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent.go:43:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent_test.go:75:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent_test.go:113:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logevent_test.go:122:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:108:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config_test.go:10:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\input_test.go:10:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\output.go:34:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config.go:113:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\output_test.go:10:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\input.go:28:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\config_test.go:26:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\input.go:33:1"
                        }
                    ]
                },
                "github.com\\wgliang\\logcool\\utils\\logo": {
                    "name": "github.com\\wgliang\\logcool\\utils\\logo",
                    "description": "1.50",
                    "errors": [
                        {
                            "line_number": 2,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logo\\logo.go:24:1"
                        },
                        {
                            "line_number": 1,
                            "error_string": "D:\\golang\\src\\github.com\\wgliang\\logcool\\utils\\logo\\logo.go:39:1"
                        }
                    ]
                }
            },
            "weight": 0.2,
            "percentage": 90,
            "error": ""
        },
        "DeadCodeTips": {
            "name": "DeadCode",
            "description": "All useless code, or never obsolete obsolete code.",
            "summaries": {},
            "weight": 0.1,
            "percentage": 90,
            "error": ""
        },
        "DependGraphTips": {
            "name": "DependGraph",
            "description": "The dependency graph for all packages in the project helps you optimize the project architecture.",
            "summaries": {
                "graph": {
                    "name": "graph",
                    "description": "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>\r\n<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\"\r\n \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\r\n<!-- Generated by graphviz version 2.38.0 (20140413.2041)\r\n -->\r\n<!-- Title: godep Pages: 1 -->\r\n<svg width=\"2423pt\" height=\"970pt\"\r\n viewBox=\"0.00 0.00 2422.80 970.00\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\">\r\n<g id=\"graph0\" class=\"graph\" transform=\"scale(1 1) rotate(0) translate(4 966)\">\r\n<title>godep</title>\r\n<polygon fill=\"white\" stroke=\"none\" points=\"-4,4 -4,-966 2418.8,-966 2418.8,4 -4,4\"/>\r\n<!-- 0 -->\r\n<g id=\"node1\" class=\"node\"><title>0</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1074.91\" cy=\"-389\" rx=\"120.779\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1074.91\" y=\"-385.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/codegangsta/inject</text>\r\n</g>\r\n<!-- 1 -->\r\n<g id=\"node2\" class=\"node\"><title>1</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-580\" rx=\"158.672\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-576.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/filter/metrics</text>\r\n</g>\r\n<!-- 2 -->\r\n<g id=\"node3\" class=\"node\"><title>2</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"751.633\" cy=\"-389\" rx=\"128.077\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"751.633\" y=\"-385.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/utils</text>\r\n</g>\r\n<!-- 1&#45;&gt;2 -->\r\n<g id=\"edge1\" class=\"edge\"><title>1&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M532.958,-569.192C546.037,-564.102 558.328,-557.225 568.746,-548 607.26,-513.9 568.514,-474.515 604.746,-438 619.221,-423.413 638.167,-413.168 657.59,-405.973\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"658.832,-409.247 667.146,-402.68 656.552,-402.629 658.832,-409.247\"/>\r\n</g>\r\n<!-- 2&#45;&gt;0 -->\r\n<g id=\"edge21\" class=\"edge\"><title>2&#45;&gt;0</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M880.071,-389C900.998,-389 922.711,-389 943.705,-389\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"943.946,-392.5 953.946,-389 943.946,-385.5 943.946,-392.5\"/>\r\n</g>\r\n<!-- 4 -->\r\n<g id=\"node5\" class=\"node\"><title>4</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1074.91\" cy=\"-319\" rx=\"110.48\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1074.91\" y=\"-315.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/Sirupsen/logrus</text>\r\n</g>\r\n<!-- 2&#45;&gt;4 -->\r\n<g id=\"edge20\" class=\"edge\"><title>2&#45;&gt;4</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M821.949,-373.897C873.765,-362.608 944.563,-347.182 997.919,-335.557\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"998.982,-338.907 1008.01,-333.358 997.492,-332.068 998.982,-338.907\"/>\r\n</g>\r\n<!-- 17 -->\r\n<g id=\"node18\" class=\"node\"><title>17</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1074.91\" cy=\"-489\" rx=\"96.3833\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1074.91\" y=\"-485.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/golang/glog</text>\r\n</g>\r\n<!-- 2&#45;&gt;17 -->\r\n<g id=\"edge22\" class=\"edge\"><title>2&#45;&gt;17</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M805.325,-405.401C862.828,-423.299 954.678,-451.888 1014.65,-470.555\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1013.77,-473.946 1024.36,-473.576 1015.85,-467.263 1013.77,-473.946\"/>\r\n</g>\r\n<!-- 3 -->\r\n<g id=\"node4\" class=\"node\"><title>3</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-742\" rx=\"163.271\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-738.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/input/collectd</text>\r\n</g>\r\n<!-- 3&#45;&gt;2 -->\r\n<g id=\"edge9\" class=\"edge\"><title>3&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M544.398,-732.439C553.353,-727.923 561.62,-722.2 568.746,-715 641.835,-641.146 548.318,-571.248 604.746,-484 626.632,-450.16 665.297,-425.799 697.324,-410.223\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"699.031,-413.288 706.592,-405.867 696.053,-406.953 699.031,-413.288\"/>\r\n</g>\r\n<!-- 3&#45;&gt;4 -->\r\n<g id=\"edge2\" class=\"edge\"><title>3&#45;&gt;4</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M542.98,-732.138C552.374,-727.716 561.121,-722.095 568.746,-715 618.624,-668.593 558.968,-616.455 604.746,-566 697.618,-463.641 800.54,-561.481 898.52,-464 932.6,-430.093 899.853,-395.306 934.52,-362 947.695,-349.343 964.513,-340.405 981.875,-334.095\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"983.186,-337.347 991.557,-330.852 980.963,-330.709 983.186,-337.347\"/>\r\n</g>\r\n<!-- 5 -->\r\n<g id=\"node6\" class=\"node\"><title>5</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1379.33\" cy=\"-668\" rx=\"124.278\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1379.33\" y=\"-664.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/shirou/gopsutil/cpu</text>\r\n</g>\r\n<!-- 3&#45;&gt;5 -->\r\n<g id=\"edge3\" class=\"edge\"><title>3&#45;&gt;5</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M515.94,-728.726C533.828,-725.154 551.982,-720.655 568.746,-715 585.889,-709.218 587.267,-700.667 604.746,-696 834.167,-634.748 1115.41,-644.709 1267.65,-656.76\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1267.69,-660.275 1277.94,-657.594 1268.26,-653.298 1267.69,-660.275\"/>\r\n</g>\r\n<!-- 6 -->\r\n<g id=\"node7\" class=\"node\"><title>6</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1379.33\" cy=\"-576\" rx=\"125.378\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1379.33\" y=\"-572.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/shirou/gopsutil/disk</text>\r\n</g>\r\n<!-- 3&#45;&gt;6 -->\r\n<g id=\"edge4\" class=\"edge\"><title>3&#45;&gt;6</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M537.473,-731.268C548.609,-727.133 559.233,-721.811 568.746,-715 596.232,-695.321 576.568,-666.674 604.746,-648 818.93,-506.053 1142.71,-535.492 1294.38,-559.85\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1294.13,-563.355 1304.56,-561.519 1295.26,-556.447 1294.13,-563.355\"/>\r\n</g>\r\n<!-- 7 -->\r\n<g id=\"node8\" class=\"node\"><title>7</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"751.633\" cy=\"-723\" rx=\"126.178\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"751.633\" y=\"-719.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/shirou/gopsutil/host</text>\r\n</g>\r\n<!-- 3&#45;&gt;7 -->\r\n<g id=\"edge5\" class=\"edge\"><title>3&#45;&gt;7</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M552.009,-733.97C575.81,-732.655 600.309,-731.302 623.621,-730.015\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"624.095,-733.494 633.886,-729.448 623.709,-726.505 624.095,-733.494\"/>\r\n</g>\r\n<!-- 8 -->\r\n<g id=\"node9\" class=\"node\"><title>8</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1379.33\" cy=\"-944\" rx=\"128.077\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1379.33\" y=\"-940.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/shirou/gopsutil/mem</text>\r\n</g>\r\n<!-- 3&#45;&gt;8 -->\r\n<g id=\"edge6\" class=\"edge\"><title>3&#45;&gt;8</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M446.469,-759.497C485.96,-776.206 548.513,-800.998 604.746,-816 835.538,-877.571 1113.39,-914.872 1265.35,-932.248\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1265.3,-935.764 1275.63,-933.414 1266.09,-928.809 1265.3,-935.764\"/>\r\n</g>\r\n<!-- 9 -->\r\n<g id=\"node10\" class=\"node\"><title>9</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1379.33\" cy=\"-890\" rx=\"120.779\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1379.33\" y=\"-886.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/shirou/gopsutil/net</text>\r\n</g>\r\n<!-- 3&#45;&gt;9 -->\r\n<g id=\"edge7\" class=\"edge\"><title>3&#45;&gt;9</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M484.562,-757.85C520.998,-765.028 565.016,-773.369 604.746,-780 844.925,-820.084 1128.62,-858.037 1276.72,-877.141\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1276.35,-880.621 1286.72,-878.427 1277.24,-873.679 1276.35,-880.621\"/>\r\n</g>\r\n<!-- 10 -->\r\n<g id=\"node11\" class=\"node\"><title>10</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1074.91\" cy=\"-789\" rx=\"139.175\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1074.91\" y=\"-785.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/shirou/gopsutil/process</text>\r\n</g>\r\n<!-- 3&#45;&gt;10 -->\r\n<g id=\"edge8\" class=\"edge\"><title>3&#45;&gt;10</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M543.662,-751.653C660.381,-759.874 826.995,-771.609 942.167,-779.721\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"942.23,-783.234 952.451,-780.445 942.722,-776.251 942.23,-783.234\"/>\r\n</g>\r\n<!-- 24 -->\r\n<g id=\"node25\" class=\"node\"><title>24</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1714.31\" cy=\"-645\" rx=\"127.277\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1714.31\" y=\"-641.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/StackExchange/wmi</text>\r\n</g>\r\n<!-- 5&#45;&gt;24 -->\r\n<g id=\"edge24\" class=\"edge\"><title>5&#45;&gt;24</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1491.67,-660.309C1522.98,-658.147 1557.21,-655.782 1589.21,-653.572\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1589.78,-657.041 1599.52,-652.86 1589.3,-650.058 1589.78,-657.041\"/>\r\n</g>\r\n<!-- 25 -->\r\n<g id=\"node26\" class=\"node\"><title>25</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1714.31\" cy=\"-779\" rx=\"170.87\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1714.31\" y=\"-775.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/shirou/gopsutil/internal/common</text>\r\n</g>\r\n<!-- 5&#45;&gt;25 -->\r\n<g id=\"edge25\" class=\"edge\"><title>5&#45;&gt;25</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1446.69,-683.124C1466.36,-688.034 1487.86,-693.823 1507.37,-700 1561.23,-717.045 1621.39,-740.676 1662.71,-757.645\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1661.44,-760.904 1672.02,-761.482 1664.11,-754.433 1661.44,-760.904\"/>\r\n</g>\r\n<!-- 6&#45;&gt;24 -->\r\n<g id=\"edge30\" class=\"edge\"><title>6&#45;&gt;24</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1451.77,-590.802C1504.56,-601.743 1576.54,-616.66 1631.63,-628.074\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1631.21,-631.561 1641.71,-630.163 1632.63,-624.706 1631.21,-631.561\"/>\r\n</g>\r\n<!-- 6&#45;&gt;25 -->\r\n<g id=\"edge31\" class=\"edge\"><title>6&#45;&gt;25</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1419.47,-593.065C1445.31,-605.031 1479.41,-622.201 1507.37,-641 1524.89,-652.781 1526.39,-659.451 1543.37,-672 1585.97,-703.479 1637.86,-735.279 1673.01,-755.933\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1671.61,-759.171 1682.01,-761.195 1675.14,-753.127 1671.61,-759.171\"/>\r\n</g>\r\n<!-- 7&#45;&gt;10 -->\r\n<g id=\"edge42\" class=\"edge\"><title>7&#45;&gt;10</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M824.334,-737.73C873.313,-747.791 938.399,-761.162 989.824,-771.727\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"989.419,-775.216 999.918,-773.8 990.827,-768.359 989.419,-775.216\"/>\r\n</g>\r\n<!-- 7&#45;&gt;24 -->\r\n<g id=\"edge40\" class=\"edge\"><title>7&#45;&gt;24</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M821.569,-708.004C917.387,-688.001 1096.52,-653.641 1251.3,-641 1364.73,-631.736 1393.57,-639.895 1507.37,-641 1530.29,-641.223 1554.64,-641.584 1578.27,-642.002\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1578.29,-645.503 1588.35,-642.184 1578.41,-638.504 1578.29,-645.503\"/>\r\n</g>\r\n<!-- 7&#45;&gt;25 -->\r\n<g id=\"edge41\" class=\"edge\"><title>7&#45;&gt;25</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M868.513,-729.755C1039.74,-739.736 1360.86,-758.455 1554.08,-769.718\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1554.14,-773.228 1564.33,-770.316 1554.55,-766.239 1554.14,-773.228\"/>\r\n</g>\r\n<!-- 8&#45;&gt;25 -->\r\n<g id=\"edge26\" class=\"edge\"><title>8&#45;&gt;25</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1465.43,-930.615C1479.72,-927.045 1494.18,-922.574 1507.37,-917 1576.11,-887.948 1646.2,-834.945 1684.28,-803.762\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1686.89,-806.152 1692.36,-797.081 1682.43,-800.758 1686.89,-806.152\"/>\r\n</g>\r\n<!-- 9&#45;&gt;25 -->\r\n<g id=\"edge32\" class=\"edge\"><title>9&#45;&gt;25</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1446.36,-874.958C1466.12,-870.03 1487.75,-864.211 1507.37,-858 1561.23,-840.955 1621.39,-817.324 1662.71,-800.355\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1664.11,-803.567 1672.02,-796.518 1661.44,-797.096 1664.11,-803.567\"/>\r\n</g>\r\n<!-- 10&#45;&gt;5 -->\r\n<g id=\"edge44\" class=\"edge\"><title>10&#45;&gt;5</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1118.92,-771.792C1173.6,-749.913 1267.93,-712.171 1326.52,-688.73\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1327.98,-691.918 1335.96,-684.953 1325.38,-685.419 1327.98,-691.918\"/>\r\n</g>\r\n<!-- 10&#45;&gt;8 -->\r\n<g id=\"edge46\" class=\"edge\"><title>10&#45;&gt;8</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1094.69,-806.969C1125.06,-834.968 1187.96,-888.722 1251.3,-917 1261.07,-921.364 1271.59,-925.037 1282.23,-928.127\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1281.46,-931.545 1292.03,-930.803 1283.31,-924.792 1281.46,-931.545\"/>\r\n</g>\r\n<!-- 10&#45;&gt;9 -->\r\n<g id=\"edge47\" class=\"edge\"><title>10&#45;&gt;9</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1115.43,-806.233C1150.64,-821.176 1203.79,-842.74 1251.3,-858 1267.59,-863.233 1285.28,-868.154 1302.13,-872.494\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1301.39,-875.917 1311.94,-874.981 1303.11,-869.131 1301.39,-875.917\"/>\r\n</g>\r\n<!-- 10&#45;&gt;24 -->\r\n<g id=\"edge43\" class=\"edge\"><title>10&#45;&gt;24</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1148.55,-773.723C1234.39,-755.541 1381.49,-723.99 1507.37,-695 1551.66,-684.8 1601.16,-672.802 1640.53,-663.115\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1641.44,-666.496 1650.31,-660.705 1639.76,-659.7 1641.44,-666.496\"/>\r\n</g>\r\n<!-- 10&#45;&gt;25 -->\r\n<g id=\"edge45\" class=\"edge\"><title>10&#45;&gt;25</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1138.34,-805.02C1171.62,-812.744 1213.41,-821.152 1251.3,-825 1364.52,-836.501 1394.25,-837.481 1507.37,-825 1553.53,-819.908 1604.47,-808.58 1644.13,-798.441\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1645.16,-801.79 1653.97,-795.893 1643.41,-795.014 1645.16,-801.79\"/>\r\n</g>\r\n<!-- 26 -->\r\n<g id=\"node27\" class=\"node\"><title>26</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1379.33\" cy=\"-798\" rx=\"94.4839\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1379.33\" y=\"-794.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/shirou/w32</text>\r\n</g>\r\n<!-- 10&#45;&gt;26 -->\r\n<g id=\"edge48\" class=\"edge\"><title>10&#45;&gt;26</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1210.76,-793.013C1232.55,-793.662 1254.8,-794.324 1275.68,-794.945\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1275.7,-798.447 1285.8,-795.246 1275.91,-791.45 1275.7,-798.447\"/>\r\n</g>\r\n<!-- 11 -->\r\n<g id=\"node12\" class=\"node\"><title>11</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"2312.11\" cy=\"-645\" rx=\"102.882\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"2312.11\" y=\"-641.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/go&#45;ole/go&#45;ole</text>\r\n</g>\r\n<!-- 12 -->\r\n<g id=\"node13\" class=\"node\"><title>12</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-218\" rx=\"129.977\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-214.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/cmd</text>\r\n</g>\r\n<!-- 12&#45;&gt;2 -->\r\n<g id=\"edge10\" class=\"edge\"><title>12&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M518.764,-226.906C536.126,-231.017 553.377,-236.828 568.746,-245 589.287,-255.922 587.262,-267.653 604.746,-283 640.188,-314.109 684.431,-345.308 714.846,-365.72\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"713.165,-368.806 723.427,-371.439 717.047,-362.981 713.165,-368.806\"/>\r\n</g>\r\n<!-- 13 -->\r\n<g id=\"node14\" class=\"node\"><title>13</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"751.633\" cy=\"-256\" rx=\"146.774\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"751.633\" y=\"-252.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/utils/logo</text>\r\n</g>\r\n<!-- 12&#45;&gt;13 -->\r\n<g id=\"edge11\" class=\"edge\"><title>12&#45;&gt;13</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M507.968,-229.196C546.795,-233.484 591.435,-238.415 631.756,-242.869\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"631.584,-246.371 641.908,-243.99 632.353,-239.414 631.584,-246.371\"/>\r\n</g>\r\n<!-- 14 -->\r\n<g id=\"node15\" class=\"node\"><title>14</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-164\" rx=\"157.872\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-160.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/output/email</text>\r\n</g>\r\n<!-- 14&#45;&gt;2 -->\r\n<g id=\"edge12\" class=\"edge\"><title>14&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M538.774,-173.711C549.576,-178.068 559.77,-183.72 568.746,-191 602.849,-218.657 577.905,-248.252 604.746,-283 632.474,-318.896 675.542,-348.112 707.707,-366.712\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"706.042,-369.792 716.468,-371.674 709.491,-363.7 706.042,-369.792\"/>\r\n</g>\r\n<!-- 15 -->\r\n<g id=\"node16\" class=\"node\"><title>15</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"751.633\" cy=\"-202\" rx=\"79.8859\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"751.633\" y=\"-198.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">gopkg.in/gomail.v2</text>\r\n</g>\r\n<!-- 14&#45;&gt;15 -->\r\n<g id=\"edge13\" class=\"edge\"><title>14&#45;&gt;15</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M519.733,-176.495C568.607,-181.894 624.694,-188.089 669.35,-193.022\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"669.245,-196.531 679.569,-194.15 670.014,-189.574 669.245,-196.531\"/>\r\n</g>\r\n<!-- 16 -->\r\n<g id=\"node17\" class=\"node\"><title>16</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-688\" rx=\"148.374\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-684.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/input/http</text>\r\n</g>\r\n<!-- 16&#45;&gt;2 -->\r\n<g id=\"edge16\" class=\"edge\"><title>16&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M540.14,-680.405C550.726,-675.623 560.494,-669.294 568.746,-661 639.559,-589.835 539.809,-514.565 604.746,-438 617.329,-423.164 634.555,-412.804 652.765,-405.573\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"654.165,-408.788 662.358,-402.07 651.764,-402.213 654.165,-408.788\"/>\r\n</g>\r\n<!-- 16&#45;&gt;4 -->\r\n<g id=\"edge14\" class=\"edge\"><title>16&#45;&gt;4</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M539.397,-680.219C550.205,-675.49 560.224,-669.223 568.746,-661 626.513,-605.256 548.748,-541.52 604.746,-484 698.233,-387.973 788.3,-492.239 898.52,-416 922.243,-399.591 911.402,-379.249 934.52,-362 950.582,-350.016 970.032,-341.332 989.241,-335.054\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"990.408,-338.356 998.939,-332.073 988.352,-331.665 990.408,-338.356\"/>\r\n</g>\r\n<!-- 16&#45;&gt;17 -->\r\n<g id=\"edge15\" class=\"edge\"><title>16&#45;&gt;17</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M536.451,-679.343C548.129,-674.846 559.147,-668.871 568.746,-661 603.662,-632.371 569.797,-594.588 604.746,-566 711.713,-478.502 879.436,-472.517 982.694,-478.851\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"982.58,-482.351 992.791,-479.522 983.045,-475.366 982.58,-482.351\"/>\r\n</g>\r\n<!-- 18 -->\r\n<g id=\"node19\" class=\"node\"><title>18</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"2047.33\" cy=\"-622\" rx=\"126.178\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"2047.33\" y=\"-618.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/go&#45;ole/go&#45;ole/oleutil</text>\r\n</g>\r\n<!-- 18&#45;&gt;11 -->\r\n<g id=\"edge17\" class=\"edge\"><title>18&#45;&gt;11</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M2155.21,-631.355C2173.15,-632.925 2191.72,-634.55 2209.52,-636.108\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"2209.63,-639.631 2219.89,-637.016 2210.24,-632.658 2209.63,-639.631\"/>\r\n</g>\r\n<!-- 19 -->\r\n<g id=\"node20\" class=\"node\"><title>19</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-488\" rx=\"157.872\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-484.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/output/lexec</text>\r\n</g>\r\n<!-- 19&#45;&gt;2 -->\r\n<g id=\"edge18\" class=\"edge\"><title>19&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M518.959,-475.367C535.945,-471.726 553.03,-467.038 568.746,-461 586.47,-454.191 587.619,-446.195 604.746,-438 629.321,-426.242 657.464,-416.076 682.433,-408.101\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"683.698,-411.372 692.193,-405.042 681.604,-404.693 683.698,-411.372\"/>\r\n</g>\r\n<!-- 20 -->\r\n<g id=\"node21\" class=\"node\"><title>20</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"751.633\" cy=\"-18\" rx=\"131.877\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"751.633\" y=\"-14.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/garyburd/redigo/redis</text>\r\n</g>\r\n<!-- 21 -->\r\n<g id=\"node22\" class=\"node\"><title>21</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"1074.91\" cy=\"-18\" rx=\"140.275\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"1074.91\" y=\"-14.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/garyburd/redigo/internal</text>\r\n</g>\r\n<!-- 20&#45;&gt;21 -->\r\n<g id=\"edge19\" class=\"edge\"><title>20&#45;&gt;21</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M883.771,-18C897.045,-18 910.596,-18 924.047,-18\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"924.306,-21.5001 934.306,-18 924.306,-14.5001 924.306,-21.5001\"/>\r\n</g>\r\n<!-- 22 -->\r\n<g id=\"node23\" class=\"node\"><title>22</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"751.633\" cy=\"-148\" rx=\"112.38\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"751.633\" y=\"-144.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/fsnotify/fsnotify</text>\r\n</g>\r\n<!-- 23 -->\r\n<g id=\"node24\" class=\"node\"><title>23</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-434\" rx=\"148.374\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-430.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/filter/zeus</text>\r\n</g>\r\n<!-- 23&#45;&gt;2 -->\r\n<g id=\"edge23\" class=\"edge\"><title>23&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M507.041,-420.863C551.108,-415.099 602.804,-408.337 647.208,-402.529\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"647.69,-405.995 657.151,-401.228 646.782,-399.055 647.69,-405.995\"/>\r\n</g>\r\n<!-- 24&#45;&gt;11 -->\r\n<g id=\"edge28\" class=\"edge\"><title>24&#45;&gt;11</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1840.27,-647.816C1866.99,-648.309 1895.05,-648.746 1921.24,-649 2033.32,-650.088 2061.35,-650.627 2173.42,-649 2182.38,-648.87 2191.66,-648.692 2200.98,-648.484\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"2201.18,-651.98 2211.1,-648.246 2201.02,-644.982 2201.18,-651.98\"/>\r\n</g>\r\n<!-- 24&#45;&gt;18 -->\r\n<g id=\"edge29\" class=\"edge\"><title>24&#45;&gt;18</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M1829.23,-637.085C1859.56,-634.977 1892.45,-632.692 1923.25,-630.552\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1923.87,-634.017 1933.6,-629.833 1923.38,-627.034 1923.87,-634.017\"/>\r\n</g>\r\n<!-- 27 -->\r\n<g id=\"node28\" class=\"node\"><title>27</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-380\" rx=\"162.471\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-376.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/output/stdout</text>\r\n</g>\r\n<!-- 27&#45;&gt;2 -->\r\n<g id=\"edge27\" class=\"edge\"><title>27&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M564.023,-384.118C581.092,-384.565 598.355,-385.016 615.131,-385.455\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"615.42,-388.964 625.508,-385.727 615.603,-381.966 615.42,-388.964\"/>\r\n</g>\r\n<!-- 28 -->\r\n<g id=\"node29\" class=\"node\"><title>28</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-110\" rx=\"144.874\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-106.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/input/file</text>\r\n</g>\r\n<!-- 28&#45;&gt;2 -->\r\n<g id=\"edge35\" class=\"edge\"><title>28&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M537.466,-117.492C548.945,-122.218 559.639,-128.566 568.746,-137 617.782,-182.41 567.447,-227.544 604.746,-283 629.661,-320.043 672.083,-348.707 704.779,-366.833\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"703.258,-369.99 713.72,-371.664 706.586,-363.831 703.258,-369.99\"/>\r\n</g>\r\n<!-- 28&#45;&gt;4 -->\r\n<g id=\"edge33\" class=\"edge\"><title>28&#45;&gt;4</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M533.495,-101.486C657.379,-95.1875 836.085,-92.0541 898.52,-121 976.7,-157.246 1034.5,-246.742 1059.8,-291.947\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1056.82,-293.798 1064.7,-300.885 1062.96,-290.434 1056.82,-293.798\"/>\r\n</g>\r\n<!-- 28&#45;&gt;22 -->\r\n<g id=\"edge34\" class=\"edge\"><title>28&#45;&gt;22</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M514.518,-121.919C557.388,-126.655 606.458,-132.075 648.718,-136.743\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"648.42,-140.231 658.743,-137.85 649.188,-133.273 648.42,-140.231\"/>\r\n</g>\r\n<!-- 29 -->\r\n<g id=\"node30\" class=\"node\"><title>29</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-272\" rx=\"151.373\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-268.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/input/stdin</text>\r\n</g>\r\n<!-- 29&#45;&gt;2 -->\r\n<g id=\"edge37\" class=\"edge\"><title>29&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M519.706,-283.938C536.508,-287.672 553.34,-292.572 568.746,-299 586.961,-306.599 587.613,-315.203 604.746,-325 633.848,-341.641 668.096,-356.871 696.094,-368.296\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"695.134,-371.683 705.717,-372.175 697.751,-365.19 695.134,-371.683\"/>\r\n</g>\r\n<!-- 29&#45;&gt;4 -->\r\n<g id=\"edge36\" class=\"edge\"><title>29&#45;&gt;4</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M536.46,-281.146C660.615,-289.89 845.798,-302.933 963.162,-311.2\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"963.027,-314.699 973.248,-311.91 963.519,-307.716 963.027,-314.699\"/>\r\n</g>\r\n<!-- 30 -->\r\n<g id=\"node31\" class=\"node\"><title>30</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-326\" rx=\"147.574\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-322.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/filter/split</text>\r\n</g>\r\n<!-- 30&#45;&gt;2 -->\r\n<g id=\"edge38\" class=\"edge\"><title>30&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M495.473,-340.298C519.297,-344.268 545.024,-348.678 568.746,-353 601.323,-358.935 637.088,-365.913 668.218,-372.137\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"667.66,-375.595 678.153,-374.13 669.036,-368.732 667.66,-375.595\"/>\r\n</g>\r\n<!-- 31 -->\r\n<g id=\"node32\" class=\"node\"><title>31</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-634\" rx=\"148.674\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-630.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/filter/grok</text>\r\n</g>\r\n<!-- 31&#45;&gt;2 -->\r\n<g id=\"edge39\" class=\"edge\"><title>31&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M539.708,-625.965C550.383,-621.273 560.29,-615.085 568.746,-607 624.254,-553.929 553.629,-495.313 604.746,-438 617.982,-423.16 635.853,-412.807 654.544,-405.587\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"655.827,-408.845 664.07,-402.189 653.475,-402.252 655.827,-408.845\"/>\r\n</g>\r\n<!-- 32 -->\r\n<g id=\"node33\" class=\"node\"><title>32</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"405.611\" cy=\"-56\" rx=\"157.072\" ry=\"18\"/>\r\n<text text-anchor=\"middle\" x=\"405.611\" y=\"-52.3\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">github.com/wgliang/logcool/output/redis</text>\r\n</g>\r\n<!-- 32&#45;&gt;2 -->\r\n<g id=\"edge51\" class=\"edge\"><title>32&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M542.684,-65.0032C552.266,-69.5886 561.13,-75.4829 568.746,-83 633.027,-146.445 556.75,-206.491 604.746,-283 628.318,-320.575 670.483,-349.026 703.441,-366.943\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"701.993,-370.137 712.469,-371.716 705.264,-363.948 701.993,-370.137\"/>\r\n</g>\r\n<!-- 32&#45;&gt;4 -->\r\n<g id=\"edge49\" class=\"edge\"><title>32&#45;&gt;4</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M562.204,-54.1559C684.102,-55.1879 842.172,-62.7304 898.52,-93 982.196,-137.949 1038.73,-241.872 1061.98,-291.566\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"1058.84,-293.113 1066.19,-300.749 1065.21,-290.199 1058.84,-293.113\"/>\r\n</g>\r\n<!-- 32&#45;&gt;20 -->\r\n<g id=\"edge50\" class=\"edge\"><title>32&#45;&gt;20</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M519.733,-43.5047C557.845,-39.295 600.343,-34.6007 638.388,-30.3984\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"638.918,-33.8612 648.474,-29.2843 638.15,-26.9035 638.918,-33.8612\"/>\r\n</g>\r\n<!-- 33 -->\r\n<g id=\"node34\" class=\"node\"><title>33</title>\r\n<ellipse fill=\"paleturquoise\" stroke=\"paleturquoise\" cx=\"103.238\" cy=\"-407\" rx=\"103.476\" ry=\"26.7407\"/>\r\n<text text-anchor=\"start\" x=\"38.2376\" y=\"-410.8\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">..srcgithub.comwgliang</text>\r\n<text text-anchor=\"middle\" x=\"103.238\" y=\"-395.8\" font-family=\"Times New Roman,serif\" font-size=\"14.00\">ogcool</text>\r\n</g>\r\n<!-- 33&#45;&gt;1 -->\r\n<g id=\"edge54\" class=\"edge\"><title>33&#45;&gt;1</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M122.362,-433.532C146.023,-465.966 190.591,-519.986 242.475,-548 255.193,-554.867 269.184,-560.27 283.454,-564.52\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"282.794,-567.969 293.366,-567.282 284.673,-561.226 282.794,-567.969\"/>\r\n</g>\r\n<!-- 33&#45;&gt;2 -->\r\n<g id=\"edge65\" class=\"edge\"><title>33&#45;&gt;2</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M127.596,-433.493C152.859,-459.955 195.825,-498.937 242.475,-515 311.03,-538.605 507.577,-553.926 568.746,-515 600.618,-494.718 576.564,-463.158 604.746,-438 620.835,-423.638 641.281,-413.462 661.738,-406.259\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"663.042,-409.515 671.44,-403.056 660.847,-402.867 663.042,-409.515\"/>\r\n</g>\r\n<!-- 33&#45;&gt;3 -->\r\n<g id=\"edge57\" class=\"edge\"><title>33&#45;&gt;3</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M107.791,-434.032C117.332,-494.952 149.421,-643.231 242.475,-715 248.751,-719.84 255.573,-723.948 262.75,-727.425\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"261.452,-730.678 272.015,-731.5 264.271,-724.27 261.452,-730.678\"/>\r\n</g>\r\n<!-- 33&#45;&gt;12 -->\r\n<g id=\"edge52\" class=\"edge\"><title>33&#45;&gt;12</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M118.78,-380.381C140.485,-343.432 185.029,-277.488 242.475,-245 254.368,-238.274 267.483,-233.12 280.948,-229.184\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"282.057,-232.51 290.802,-226.529 280.236,-225.751 282.057,-232.51\"/>\r\n</g>\r\n<!-- 33&#45;&gt;14 -->\r\n<g id=\"edge61\" class=\"edge\"><title>33&#45;&gt;14</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M113.12,-380.253C130.142,-333.497 172.054,-237.466 242.475,-191 250.342,-185.809 258.877,-181.509 267.778,-177.955\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"268.986,-181.24 277.173,-174.515 266.579,-174.666 268.986,-181.24\"/>\r\n</g>\r\n<!-- 33&#45;&gt;16 -->\r\n<g id=\"edge59\" class=\"edge\"><title>33&#45;&gt;16</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M110.566,-434.128C124.458,-487.169 162.927,-604.366 242.475,-661 250.412,-666.65 259.115,-671.274 268.239,-675.044\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"267.294,-678.428 277.886,-678.671 269.758,-671.875 267.294,-678.428\"/>\r\n</g>\r\n<!-- 33&#45;&gt;19 -->\r\n<g id=\"edge62\" class=\"edge\"><title>33&#45;&gt;19</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M155.886,-430.121C181.502,-440.825 213.163,-452.928 242.475,-461 258.459,-465.401 275.539,-469.196 292.393,-472.432\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"291.809,-475.883 302.281,-474.273 293.091,-469.001 291.809,-475.883\"/>\r\n</g>\r\n<!-- 33&#45;&gt;23 -->\r\n<g id=\"edge56\" class=\"edge\"><title>33&#45;&gt;23</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M200.931,-415.692C224.853,-417.842 250.857,-420.179 276.034,-422.442\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"275.737,-425.93 286.01,-423.339 276.363,-418.958 275.737,-425.93\"/>\r\n</g>\r\n<!-- 33&#45;&gt;27 -->\r\n<g id=\"edge64\" class=\"edge\"><title>33&#45;&gt;27</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M200.931,-398.308C222.593,-396.361 245.963,-394.261 268.88,-392.201\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"269.37,-395.671 279.017,-391.289 268.744,-388.699 269.37,-395.671\"/>\r\n</g>\r\n<!-- 33&#45;&gt;28 -->\r\n<g id=\"edge58\" class=\"edge\"><title>33&#45;&gt;28</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M109.631,-379.921C122.198,-324.467 158.94,-198.074 242.475,-137 250.595,-131.063 259.553,-126.249 268.965,-122.362\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"270.382,-125.568 278.513,-118.776 267.921,-119.015 270.382,-125.568\"/>\r\n</g>\r\n<!-- 33&#45;&gt;29 -->\r\n<g id=\"edge60\" class=\"edge\"><title>33&#45;&gt;29</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M129.644,-380.626C155.554,-355.428 198.144,-318.485 242.475,-299 254.341,-293.785 267.101,-289.57 280.063,-286.165\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"281.131,-289.507 290.004,-283.716 279.457,-282.71 281.131,-289.507\"/>\r\n</g>\r\n<!-- 33&#45;&gt;30 -->\r\n<g id=\"edge55\" class=\"edge\"><title>33&#45;&gt;30</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M155.886,-383.879C181.502,-373.175 213.163,-361.072 242.475,-353 259.149,-348.409 277.015,-344.478 294.575,-341.152\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"295.662,-344.511 304.867,-339.264 294.399,-337.626 295.662,-344.511\"/>\r\n</g>\r\n<!-- 33&#45;&gt;31 -->\r\n<g id=\"edge53\" class=\"edge\"><title>33&#45;&gt;31</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M114.575,-433.892C133.019,-477.95 176.011,-564.818 242.475,-607 251.33,-612.62 260.985,-617.183 271.027,-620.876\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"270.166,-624.279 280.76,-624.157 272.402,-617.646 270.166,-624.279\"/>\r\n</g>\r\n<!-- 33&#45;&gt;32 -->\r\n<g id=\"edge63\" class=\"edge\"><title>33&#45;&gt;32</title>\r\n<path fill=\"none\" stroke=\"black\" d=\"M107.132,-380.025C115.435,-316.94 145.36,-159.351 242.475,-83 248.942,-77.9156 256.009,-73.6339 263.458,-70.0387\"/>\r\n<polygon fill=\"black\" stroke=\"black\" points=\"265.314,-73.0478 273.082,-65.843 262.516,-66.6311 265.314,-73.0478\"/>\r\n</g>\r\n</g>\r\n</svg>\r\n",
                    "errors": null
                }
            },
            "weight": 0,
            "percentage": 90,
            "error": ""
        },
        "ImportPackagesTips": {
            "name": "ImportPackages",
            "description": "Check the project variables, functions, etc. naming spelling is wrong.",
            "summaries": {},
            "weight": 0,
            "percentage": 90,
            "error": ""
        },
        "SimpleTips": {
            "name": "Simple",
            "description": "All golang code hints that can be optimized and give suggestions for changes.",
            "summaries": {},
            "weight": 0.1,
            "percentage": 90,
            "error": ""
        },
        "SpellCheckTips": {
            "name": "SpellCheck",
            "description": "Check the project variables, functions, etc. naming spelling is wrong.",
            "summaries": {},
            "weight": 0.1,
            "percentage": 90,
            "error": ""
        },
        "UnitTestTips": {
            "name": "UnitTest",
            "description": "Run all valid TEST in your golang package.And will measure from both coverage and time-consuming.",
            "summaries": {
                "github.com\\wgliang\\logcool\\cmd": {
                    "name": "github.com\\wgliang\\logcool\\cmd",
                    "description": "{\"is_pass\":true,\"coverage\":\"58.0%\",\"time\":5.459}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\filter\\grok": {
                    "name": "github.com\\wgliang\\logcool\\filter\\grok",
                    "description": "{\"is_pass\":true,\"coverage\":\"80.0%\",\"time\":8.581}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\filter\\metrics": {
                    "name": "github.com\\wgliang\\logcool\\filter\\metrics",
                    "description": "{\"is_pass\":true,\"coverage\":\"88.2%\",\"time\":7.092}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\filter\\split": {
                    "name": "github.com\\wgliang\\logcool\\filter\\split",
                    "description": "{\"is_pass\":true,\"coverage\":\"92.3%\",\"time\":3.753}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\filter\\zeus": {
                    "name": "github.com\\wgliang\\logcool\\filter\\zeus",
                    "description": "{\"is_pass\":true,\"coverage\":\"83.3%\",\"time\":9.903}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\input\\collectd": {
                    "name": "github.com\\wgliang\\logcool\\input\\collectd",
                    "description": "{\"is_pass\":true,\"coverage\":\"36.4%\",\"time\":8.873}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\input\\file": {
                    "name": "github.com\\wgliang\\logcool\\input\\file",
                    "description": "{\"is_pass\":true,\"coverage\":\"9.5%\",\"time\":10.097}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\input\\http": {
                    "name": "github.com\\wgliang\\logcool\\input\\http",
                    "description": "{\"is_pass\":true,\"coverage\":\"81.1%\",\"time\":15.025}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\input\\stdin": {
                    "name": "github.com\\wgliang\\logcool\\input\\stdin",
                    "description": "{\"is_pass\":true,\"coverage\":\"75.0%\",\"time\":20.609}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\output\\email": {
                    "name": "github.com\\wgliang\\logcool\\output\\email",
                    "description": "{\"is_pass\":true,\"coverage\":\"84.0%\",\"time\":19.4}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\output\\lexec": {
                    "name": "github.com\\wgliang\\logcool\\output\\lexec",
                    "description": "{\"is_pass\":true,\"coverage\":\"73.3%\",\"time\":9.581}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\output\\redis": {
                    "name": "github.com\\wgliang\\logcool\\output\\redis",
                    "description": "{\"is_pass\":true,\"coverage\":\"28.2%\",\"time\":11.354}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\output\\stdout": {
                    "name": "github.com\\wgliang\\logcool\\output\\stdout",
                    "description": "{\"is_pass\":true,\"coverage\":\"81.8%\",\"time\":17.312}",
                    "errors": null
                },
                "github.com\\wgliang\\logcool\\utils": {
                    "name": "github.com\\wgliang\\logcool\\utils",
                    "description": "{\"is_pass\":false,\"coverage\":\"0%\",\"time\":0}",
                    "errors": null
                }
            },
            "weight": 0.4,
            "percentage": 62.221428571428575,
            "error": ""
        }
    },
    "issues": 21,
    "time_stamp": "2017-05-14 22:25:12"
}`

func Test_Json2Html(t *testing.T) {
	htmldata, _ := Json2Html([]byte(structJson))
	SaveAsHtml(htmldata, `github.com\wgliang\logcool`, `..\bin`, "221242124214", defaultTpl)
}
