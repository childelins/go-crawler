<?php

class JsonRPC
{
    private $conn;

    function __construct($host, $port)
    {
        $this->conn = fsockopen($host, $port, $errno, $errstr, 3);
        if (!$this->conn) {
            return;
        }
    }

    public function Call($method, $params)
    {
        if (!$this->conn) {
            return false;
        }

        $err = fwrite($this->conn, json_encode(array(
                'method' => $method,
                'params' => array($params),
                'id' => 0,
            )) . "\n");
        if ($err === false) {
            return false;
        }

        stream_set_timeout($this->conn, 0, 3000);
        $line = fgets($this->conn);
        if ($line === false) {
            return NULL;
        }

        return json_decode($line, true);
    }
}

$client = new JsonRPC("127.0.0.1", 6001);
$args = "jsonrpc";
$r = $client->Call("App.Hi", $args);
printf("result: %s\n", $r['result']);