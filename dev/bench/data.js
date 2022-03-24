window.BENCHMARK_DATA = {
  "lastUpdate": 1648134447315,
  "repoUrl": "https://github.com/CGI-FR/SIGO",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "youen.peron@cgi.com",
            "name": "Linux User",
            "username": "youen"
          },
          "committer": {
            "email": "youen.peron@cgi.com",
            "name": "Linux User",
            "username": "youen"
          },
          "distinct": true,
          "id": "c8a92ce31904c04f62d70c759618616473fb7e94",
          "message": "perf(bench): add bench test",
          "timestamp": "2022-03-24T14:48:11Z",
          "tree_id": "a94b4541d1e88b90fcc35848f5c5b00ca438d7f8",
          "url": "https://github.com/CGI-FR/SIGO/commit/c8a92ce31904c04f62d70c759618616473fb7e94"
        },
        "date": 1648134446665,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSimpleClustering",
            "value": 152233,
            "unit": "ns/op\t   39461 B/op\t    1041 allocs/op",
            "extra": "80137 times\n2 procs"
          },
          {
            "name": "BenchmarkLongClustering",
            "value": 78478876220,
            "unit": "ns/op\t6372754064 B/op\t476824006 allocs/op",
            "extra": "1 times\n2 procs"
          }
        ]
      }
    ]
  }
}