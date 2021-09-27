<?php

namespace SwagTestEnvironment\Content\MailCatcher;

use Shopware\Core\Content\Mail\Service\MailerTransportFactory;
use Shopware\Core\Framework\DataAbstractionLayer\EntityRepositoryInterface;
use Shopware\Core\System\SystemConfig\SystemConfigService;
use Symfony\Component\Mailer\Transport\TransportInterface;

class MailTransportFactory extends MailerTransportFactory
{
    private EntityRepositoryInterface $repo;

    public function __construct(EntityRepositoryInterface $repo)
    {
        $this->repo = $repo;
    }

    public function create(?SystemConfigService $configService = null): TransportInterface
    {
        return new DatabaseTransport($this->repo);
    }

    public function fromString(string $dsn): TransportInterface
    {
        return $this->create();
    }


}
