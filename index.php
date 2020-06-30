<?php
ini_set('display_errors', '0');

require 'vendor/autoload.php';

$app = new \SBP\DemoServer\Application();
$app->run();
