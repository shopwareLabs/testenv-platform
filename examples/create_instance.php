<?php

$ch = curl_init('http://localhost/environments');
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'POST');
curl_setopt($ch, CURLOPT_HTTPHEADER, array(
    'content-type: application/json',
));
curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode([
    'installVersion' => '6.4.3.1',
    'plugin' => base64_encode(file_get_contents('https://github.com/FriendsOfShopware/FroshDevelopmentHelper/releases/download/0.3.2/FroshDevelopmentHelper.zip'))
]));
curl_exec($ch);
