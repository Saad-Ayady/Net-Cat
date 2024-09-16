package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "sync"
    "time"
)

const PortAsatHH = "2525";

var (
    clients       = make(map[net.Conn]string);
    messages      []string;
    clientMutex   sync.Mutex;
    maxCon = 10;
)

func main() {
    port := PortAsatHH;

    listener, err := net.Listen("tcp", ":"+port);
    if err != nil {
        fmt.Println("Error starting server:", err);
        os.Exit(1);
    }
    defer listener.Close();
    fmt.Println("ROOM Is Raning...")
    fmt.Printf("Listening on port :%s\n", port);

    for {
        if len(clients) >= maxCon {
            fmt.Println("Max connections Is 10 :( ");
            conn, _ := listener.Accept();
            conn.Close();
            continue;
        }

        conn, err := listener.Accept();
        if err != nil {
            fmt.Println("Error accepting connection:", err);
            continue;
        }

        clientMutex.Lock();
        clients[conn] = "";
        clientMutex.Unlock();

        go handleConnection(conn);
    }
}

func handleConnection(conn net.Conn) {
    defer func() {
        clientMutex.Lock();
        name := clients[conn];
        delete(clients, conn);
        clientMutex.Unlock();
        if name != "" {
            broadcast(fmt.Sprintf("[%s] has left the chat...\n", name));
        }
        conn.Close();
    }()

    conn.Write([]byte("Welcome to TCP-Chat!\n" +
        "         _nnnn_\n" +
        "        dGGGGMMb\n" +
        "       @p~qp~~qMb\n" +
        "       M|@||@) M|\n" +
        "       @,----.JM|\n" +
        "      JS^\\__/  qKL\n" +
        "     dZP        qKRb\n" +
        "    dZP          qKKb\n" +
        "   fZP            SMMb\n" +
        "   HZM            MMMM\n" +
        "   FqM            MMMM\n" +
        " __| \".        |\\dS\"qML\n" +
        " |    `.       | `' \\Zq\n" +
        "_)      \\.___.,|     .'\n" +
        "\\____   )MMMMMP|   .'\n" +
        "     `-'       `--'\n" +
        "[ENTER YOUR NAME]:"));

    name, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println("Error reading name:", err);
        return;
    }
    name = strings.TrimSpace(name);

    clientMutex.Lock();
    clients[conn] = name;
    clientMutex.Unlock();

    broadcast(fmt.Sprintf("[%s] has joined the chat...\n", name));

    if len(messages) > 0 {
        for _, msg := range messages {
            conn.Write([]byte(msg));
        }
    }

    scanner := bufio.NewScanner(conn);
    for scanner.Scan() {
        msg := scanner.Text();
        if msg != "" {
            timestamp := time.Now().Format("2006-01-02 15:04:05");
            formattedMsg := fmt.Sprintf("[%s][%s]: %s\n", timestamp, name, msg);
            messages = append(messages, formattedMsg);
            broadcast(formattedMsg);
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading from connection:", err);
    }
}

func broadcast(message string) {
    clientMutex.Lock();
    defer clientMutex.Unlock();
    for conn := range clients {
        if _, err := conn.Write([]byte(message)); err != nil {
            fmt.Println("Error broadcasting message:", err);
            conn.Close();
            delete(clients, conn);
        }
    }
}
