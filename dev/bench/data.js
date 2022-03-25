window.BENCHMARK_DATA = {
  "lastUpdate": 1648212610226,
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
      },
      {
        "commit": {
          "author": {
            "email": "89643755+giraud10@users.noreply.github.com",
            "name": "giraud10",
            "username": "giraud10"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4b13826921364b4601bcc49e9a584a0bc1755b43",
          "message": "feat: add option profiling (#26)",
          "timestamp": "2022-03-25T13:47:24+01:00",
          "tree_id": "efb28a9dfe0a648a2de1390d2034e05531be6c8d",
          "url": "https://github.com/CGI-FR/SIGO/commit/4b13826921364b4601bcc49e9a584a0bc1755b43"
        },
        "date": 1648212609400,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSimpleClustering",
            "value": 120886,
            "unit": "ns/op\t   39460 B/op\t    1041 allocs/op",
            "extra": "98649 times\n2 procs"
          },
          {
            "name": "BenchmarkLongClustering",
            "value": 59336403592,
            "unit": "ns/op\t6372753504 B/op\t476856764 allocs/op",
            "extra": "1 times\n2 procs"
          }
        ]
      }
    ]
  }
}