<?php

namespace SBP\DemoServer;

use Docker\Docker;
use SBP\DemoServer\Controller\Create;
use SBP\DemoServer\Controller\Delete;
use SBP\DemoServer\Controller\Index;

class Container extends \Slim\Container
{
    public function __construct(array $values = [])
    {
        $values['settings'] = [
            'displayErrorDetails' => true
        ];
        parent::__construct($values);
        $me = $this;

        $this['docker.client'] = function () {
            return Docker::create();
        };

        $this['docker.service'] = function () use($me) {
            return new \SBP\DemoServer\Services\Docker($me['docker.client']);
        };

        $this['controller.index'] = new Index($this);
        $this['controller.create'] = new Create($this);
        $this['controller.delete'] = new Delete($this);
    }
}