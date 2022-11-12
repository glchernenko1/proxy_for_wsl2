# proxy_for_wsl2
Automatically accessing a WSL2 distribution from your local area network (LAN)

When using a WSL 1 distribution, if your computer was set up to be accessed by your LAN, then applications run in WSL could be accessed on your LAN as well.

This isn't the default case in WSL 2. WSL 2 has a virtualized ethernet adapter with its own unique IP address. Currently, to enable this workflow you will need to go through the same steps as you would for a regular virtual machine.

This application reads ports from a json file and redirects traffic to wsl 2. For example, this can be useful if you want to connect to wsl 2 via ssh

# Install
You can download the binary file.

OR
```
go install github.com/glchernenko1/proxy_for_wsl2@v1.0.0
```
In the file properties you need to add parameter "run this program as an administrator"
![image](https://user-images.githubusercontent.com/42982650/201486929-cf8cda45-a772-45f0-a389-b6fa72562973.png)

# JSON Configurator
```JSON
{
  "ports": [
    {
      "listenport": "2222",
      "connectport": "8080"
    },
    {
      "listenport": "4000",
      "connectport": "4000"
    }
  ]
}
```

# Run application

```
.\proxy_for_wsl2.exe path_to_json
```
If you want the application to start when Windows starts. Create a task with admin rights in "Task scheduler". In the "add argument" column, add the full path to your JSON file. 

Detailed instructions on how to add a program to autorun https://www.windowscentral.com/how-create-automated-task-using-task-scheduler-windows-10
