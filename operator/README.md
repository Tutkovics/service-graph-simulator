## Example output:
```
                --- 0 ---
                name: front-end
                port: 80
                node port: 30000
                resource: 1000 (kB), 100 (mCPU)
                # of endpoints: 3
                        path: /instant
                        cpu: 100
                        delay: 10
                        path: /chain
                        cpu: 200
                        delay: 20
                                call out: back-end:80/profile
                        path: /iterate
                        cpu: 300
                        delay: 30
                                call out: back-end:80/profile
                                call out: monitor:80/index
                --- 1 ---
                name: back-end
                port: 80
                node port: 0
                resource: 2000 (kB), 200 (mCPU)
                # of endpoints: 2
                        path: /profile
                        cpu: 100
                        delay: 10
                                call out: db:80/get
                        path: /create
                        cpu: 200
                        delay: 20
                                call out: db:80/set
                --- 2 ---
                name: db
                port: 36
                node port: 0
                resource: 3000 (kB), 300 (mCPU)
                # of endpoints: 2
                        path: /get
                        cpu: 100
                        delay: 10
                        path: /set
                        cpu: 200
                        delay: 20
                --- 3 ---
                name: monitor
                port: 70
                node port: 0
                resource: 1000 (kB), 100 (mCPU)
                # of endpoints: 1
                        path: /index
                        cpu: 100
                        delay: 10
```