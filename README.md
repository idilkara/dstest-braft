This is a repository to test braft and some versions of it using DSTest.
Instead of using these repositories as submodules, I cloned them as a preference. 

- The links to the original repositories are here,

dstest: https://github.com/egeberkaygulcan/dstest.git

braft: https://github.com/baidu/braft.git

brpc: https://github.com/apache/brpc.git

**Results**

| index | benchmark | scheduler | fault condition | number of unique schedules | number of unique traces | number of tests with error or failure logged | unique errors or failures logged | number of possible liveness bugs | number of safety violations | mean execution time           |
|-------|-----------|-----------|------------------|-----------------------------|--------------------------|--------------------------------------------|---------------------------------|---------------------------------|------------------------------|--------------------------------|
| 0     | b0        | pct1      | False            | 52                          | 115                      | 1                                          | {E0824}                         | 94                              | 0                            | 0 days 00:00:17.027642740     |
| 1     | b0        | pct1      | True             | 73                          | 170                      | 0                                          | {}                              | 98                              | 0                            | 0 days 00:00:16.985119440     |
| 2     | b0        | random    | False            | 100                         | 300                      | 0                                          | {}                              | 6                               | 0                            | 0 days 00:00:16.745741770     |
| 3     | b0        | random    | True             | 100                         | 300                      | 2                                          | {E0824}                         | 13                              | 0                            | 0 days 00:00:17.191952710     |
| 4     | b0        | pct2      | False            | 67                          | 168                      | 0                                          | {}                              | 95                              | 0                            | 0 days 00:00:17.191525940     |
| 5     | b0        | pct2      | True             | 84                          | 212                      | 0                                          | {}                              | 98                              | 0                            | 0 days 00:00:17.190402410     |
| 6     | b0        | pos       | False            | 100                         | 300                      | 2                                          | {E0824}                         | 1                               | 0                            | 0 days 00:00:17.228563320     |
| 7     | b0        | pos       | True             | 100                         | 300                      | 1                                          | {E0824}                         | 12                              | 0                            | 0 days 00:00:17.727142940     |
| 8     | b0        | posc      | False            | 100                         | 300                      | 4                                          | {E0824}                         | 7                               | 0                            | 0 days 00:00:17.048977010     |
| 9     | b0        | posc      | True             | 100                         | 300                      | 3                                          | {E0824}                         | 14                              | 0                            | 0 days 00:00:17.170510300     |
| 10    | b1        | pct1      | False            | 52                          | 119                      | 0                                          | {}                              | 98                              | 0                            | 0 days 00:00:17.171506020     |
| 11    | b1        | pct1      | True             | 74                          | 166                      | 0                                          | {}                              | 97                              | 0                            | 0 days 00:00:17.136797490     |
| 12    | b1        | random    | False            | 100                         | 300                      | 9                                          | {E0824}                         | 6                               | 5                            | 0 days 00:00:16.849027620     |
| 13    | b1        | random    | True             | 100                         | 300                      | 10                                         | {E0824}                         | 10                              | 6                            | 0 days 00:00:17.187320950     |
| 14    | b1        | pct2      | False            | 72                          | 189                      | 0                                          | {}                              | 96                              | 0                            | 0 days 00:00:17.346660550     |
| 15    | b1        | pct2      | True             | 79                          | 205                      | 0                                          | {}                              | 97                              | 0                            | 0 days 00:00:17.311387160     |
| 16    | b1        | pos       | False            | 100                         | 300                      | 3                                          | {E0824}                         | 3                               | 2                            | 0 days 00:00:17.350821590     |
| 17    | b1        | pos       | True             | 100                         | 300                      | 7                                          | {E0824}                         | 12                              | 4                            | 0 days 00:00:17.335166660     |
| 18    | b1        | posc      | False            | 100                         | 300                      | 6                                          | {E0824}                         | 0                               | 2                            | 0 days 00:00:17.203054920     |
| 19    | b1        | posc      | True             | 100                         | 300                      | 4                                          | {E0824}                         | 12                              | 4                            | 0 days 00:00:17.176945260     |
| 20    | b2        | pct1      | False            | 56                          | 129                      | 0                                          | {}                              | 88                              | 0                            | 0 days 00:00:17.065691060     |
| 21    | b2        | pct1      | True             | 74                          | 165                      | 2                                          | {E0826}                         | 96                              | 0                            | 0 days 00:00:17.150621990     |
| 22    | b2        | random    | False            | 100                         | 300                      | 10                                         | {F0826, E0826}                  | 11                              | 0                            | 0 days 00:00:17.598116050     |
| 23    | b2        | random    | True             | 100                         | 300                      | 15                                         | {F0826, E0826}                  | 19                              | 0                            | 0 days 00:00:17.642008700     |
| 24    | b2        | pct2      | False            | 69                          | 192                      | 0                                          | {}                              | 94                              | 0                            | 0 days 00:00:17.547892460     |
| 25    | b2        | pct2      | True             | 89                          | 230                      | 1                                          | {E0826}                         | 95                              | 0                            | 0 days 00:00:17.445940140     |
| 26    | b2        | pos       | False            | 100                         | 300                      | 11                                         | {F0826, E0826}                  | 16                              | 0                            | 0 days 00:00:17.229746250     |
| 27    | b2        | pos       | True             | 100                         | 300                      | 7                                          | {F0826, E0826}                  | 17                              | 0                            | 0 days 00:00:17.692410040     |
| 28    | b2        | posc      | False            | 100                         | 300                      | 13                                         | {F0826, E0826}                  | 19                              | 0                            | 0 days 00:00:17.283015500     |
| 29    | b2        | posc      | True             | 100                         | 300                      | 9                                          | {F0826, E0826}                  | 16                              | 0                            | 0 days 00:00:17.732058390     |


Errors logged :
    E0824: 0 /src/braft_builder/brpc/src/bthread/mutex.cpp:497] bthread is suspended while holding1 pthread locks.
    E0826: /src/braft_builder/braft/src/braft/replicator.cpp:473] Group Counter fail, response term 3 mismatch, expect term 2
Failures logged : 
    F0826: Check failed: new_term >= _term (2 vs 3). 