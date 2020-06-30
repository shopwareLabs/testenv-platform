<?php


namespace SBP\DemoServer\Controller;


use SBP\DemoServer\Container;
use SBP\DemoServer\Services\Docker;

abstract class AbstractController
{
    /**
     * @var Container
     */
    protected $container;

    public function __construct(Container $container)
    {
        $this->container = $container;
    }

    public function get(string $id)
    {
        return $this->container->get($id);
    }

    public function getDocker(): Docker
    {
        return $this->container->get('docker.service');
    }
}