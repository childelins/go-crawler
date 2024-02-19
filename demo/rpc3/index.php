<?php

use Spiral\Goridge\RPC\RPC;
use Spiral\Goridge\SocketRelay;

include 'vendor/autoload.php';

$rpc = new RPC(new SocketRelay("127.0.0.1", 6001));
$result = $rpc->call("App.Hi", "Goridge RPC");
echo (string) $result;