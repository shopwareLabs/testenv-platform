<?php


namespace SBP\DemoServer\Controller;


use Slim\Http\Request;
use Slim\Http\Response;

class Index extends AbstractController
{
    public function index(Request $request, Response $response)
    {
        return $response->withJson($this->getDocker()->list());
    }
}