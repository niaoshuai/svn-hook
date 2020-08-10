# 构建过程
```shell script
$ make build-windows
```
## 使用方法 (Windows)
```json
{
  "description":"json svn data",
  "data": [
    {
      "svn_server":"[包含]",
      "jenkins_server":"[jenkins URL]"
    }
  ]
}

```

```shell script
svn-hook --REPOS=%REPOS% --TXN=%TXN% --LOG_PATH=E://svn.log --JENKINS_DATA=jenkins-data.json
```
