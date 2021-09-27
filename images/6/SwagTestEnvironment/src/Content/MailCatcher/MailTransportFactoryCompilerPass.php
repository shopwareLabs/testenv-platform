<?php

namespace SwagTestEnvironment\Content\MailCatcher;

use Symfony\Component\DependencyInjection\Compiler\CompilerPassInterface;
use Symfony\Component\DependencyInjection\ContainerBuilder;
use Symfony\Component\DependencyInjection\Reference;

class MailTransportFactoryCompilerPass implements CompilerPassInterface
{
    public function process(ContainerBuilder $container): void
    {
        $container->getDefinition('mailer.transport_factory')
            ->setClass(MailTransportFactory::class)
            ->setArguments([
                new Reference('sw_mail_catcher.repository')
            ]);
    }
}
