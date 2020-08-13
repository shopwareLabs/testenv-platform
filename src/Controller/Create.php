<?php

namespace SBP\DemoServer\Controller;

use SBP\DemoServer\Struct\CreateShopRequest;
use Slim\Http\Request;
use Slim\Http\Response;

class Create extends AbstractController
{
    public function create(Request $request, Response $response)
    {
        $shopRequest = (new CreateShopRequest())->fillFromRequest($request);

        if (version_compare($shopRequest->installVersion, '6.3.0', '>=')) {
            $shopRequest->installVersion = '6.3';
        } elseif (version_compare($shopRequest->installVersion, '6.2.0', '>=')) {
            $shopRequest->installVersion = '6.2';
        } elseif (version_compare($shopRequest->installVersion, '6.0.0', '>=')) {
            $shopRequest->installVersion = '6.1';
        } elseif (version_compare($shopRequest->installVersion, '5.6.0', '>=')) {
            $shopRequest->installVersion = '5.6';
        } else {
            $shopRequest->installVersion = '5.5';
        }

        $shop = $this->getDocker()->createDemoShop($shopRequest);

        return $response->withJson($shop);
    }
}
