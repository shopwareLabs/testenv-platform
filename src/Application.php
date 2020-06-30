<?php

namespace SBP\DemoServer;

use Slim\App;

class Application extends App
{
    public function __construct()
    {
        parent::__construct(new Container());

        $this->configureRoutes();
    }

    private function configureRoutes(): void
    {
        $this->get('/', [$this->getContainer()['controller.index'], 'index']);
        $this->post('/', [$this->getContainer()['controller.create'], 'create']);
        $this->delete('/', [$this->getContainer()['controller.delete'], 'delete']);
    }
}