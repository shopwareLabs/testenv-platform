<?php

namespace SwagTestEnvironment\Content\MailCatcher;

use Shopware\Core\Framework\Context;
use Shopware\Core\Framework\DataAbstractionLayer\EntityRepositoryInterface;
use Symfony\Component\Mailer\Envelope;
use Symfony\Component\Mailer\SentMessage;
use Symfony\Component\Mailer\Transport\TransportInterface;
use Symfony\Component\Mime\Email;
use Symfony\Component\Mime\RawMessage;

class DatabaseTransport implements TransportInterface
{
    private EntityRepositoryInterface $mailRepository;

    public function __construct(EntityRepositoryInterface $mailRepository)
    {
        $this->mailRepository = $mailRepository;
    }

    public function send(RawMessage $message, Envelope $envelope = null): ?SentMessage
    {
        if ($message instanceof Email) {
            $this->mailRepository->create([
                [
                    'sender' => [$message->getFrom()[0]->getAddress() => $message->getFrom()[0]->getName()],
                    'receiver' => $this->convertAddress($message->getTo()),
                    'subject' => $message->getSubject(),
                    'plainText' => nl2br($message->getTextBody()),
                    'htmlText' => $message->getHtmlBody(),
                    'eml' => $message->toString(),
                ]
            ], Context::createDefaultContext());
        }

        return new SentMessage($message, $envelope);
    }

    public function __toString(): string
    {
        return 'cody://blou';
    }

    private function convertAddress(array $addresses): array
    {
        $list = [];

        foreach ($addresses as $address) {
            $list[$address->getAddress()] = $address->getName();
        }

        return $list;
    }
}
