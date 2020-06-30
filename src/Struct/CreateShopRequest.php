<?php


namespace SBP\DemoServer\Struct;


class CreateShopRequest extends Struct
{
    /**
     * @var string
     */
    public $installVersion = '5.5.7';

    /**
     * @var string
     */
    public $plugin;
}
