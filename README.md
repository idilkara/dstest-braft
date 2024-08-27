This is a repository to test braft and some versions of it using DSTest.
Instead of using these repositories as submodules, I cloned them as a preference. 

- The links to the original repositories are here,

dstest: https://github.com/egeberkaygulcan/dstest.git

braft: https://github.com/baidu/braft.git

brpc: https://github.com/apache/brpc.git

**Results**

**BENCHMARK 0**

<img width="828" alt="image" src="https://github.com/user-attachments/assets/109684e3-0c44-4156-8335-80d333b4dc99">


**BENCHMARK 1**


<img width="828" alt="image" src="https://github.com/user-attachments/assets/7e17e839-4cc9-4829-9f69-501cf2a0b55e">


**BENCHMARK 2**


<img width="828" alt="image" src="https://github.com/user-attachments/assets/29ad44b4-5535-4bd5-9d67-b3469f601d71">


Errors logged :

    E0824: 0 /src/braft_builder/brpc/src/bthread/mutex.cpp:497] bthread is suspended while holding1 pthread locks.

    E0826: /src/braft_builder/braft/src/braft/replicator.cpp:473] Group Counter fail, response term {int} mismatch, expect term {int}

Failures logged : 

    F0826: Check failed: new_term >= _term ({int} vs {int}). 
