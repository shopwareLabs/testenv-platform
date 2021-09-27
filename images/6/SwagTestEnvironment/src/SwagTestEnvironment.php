<?php declare(strict_types=1);

namespace SwagTestEnvironment;

use Shopware\Core\Content\Product\Aggregate\ProductVisibility\ProductVisibilityDefinition;
use Shopware\Core\Defaults;
use Shopware\Core\Framework\Context;
use Shopware\Core\Framework\DataAbstractionLayer\Search\Criteria;
use Shopware\Core\Framework\DataAbstractionLayer\Search\Filter\EqualsFilter;
use Shopware\Core\Framework\DataAbstractionLayer\Search\Query\ScoreQuery;
use Shopware\Core\Framework\Plugin;
use Shopware\Core\Framework\Plugin\Context\DeactivateContext;
use Shopware\Core\Framework\Plugin\Context\InstallContext;
use Shopware\Core\Framework\Plugin\Context\UninstallContext;
use SwagTestEnvironment\Content\MailCatcher\MailTransportFactoryCompilerPass;
use Symfony\Component\DependencyInjection\ContainerBuilder;

class SwagTestEnvironment extends Plugin
{
    public function install(InstallContext $installContext): void
    {
        $c = Context::createDefaultContext();

        $salesChannelId = $this->container->get('sales_channel.repository')
            ->searchIds((new Criteria())->addFilter(new EqualsFilter('typeId', Defaults::SALES_CHANNEL_TYPE_STOREFRONT)), $c)
            ->firstId();

        $categoryCriteria = new Criteria();
        $categoryCriteria->addQuery(new ScoreQuery(new EqualsFilter('name', 'Industrial & Books'), 1000));
        $categoryCriteria->setLimit(2);

        $categoryId = $this->container->get('category.repository')
            ->searchIds($categoryCriteria, $c)
            ->firstId();

        if ($categoryId === null) {
            return;
        }

        $taxId = $this->container->get('tax.repository')
            ->searchIds(new Criteria(), $c)
            ->firstId();

        $productRepository = $this->container->get('product.repository');

        $productRepository->create(
            [
                [
                    'productNumber' => 'FREE-PRODUCT',
                    'stock' => 1,
                    'name' => 'FREE Product',
                    'active' => true,
                    'price' => [
                        [
                            'currencyId' => Defaults::CURRENCY,
                            'gross' => 0,
                            'net' => 0,
                            'linked' => true,
                        ],
                    ],
                    'manufacturer' => ['name' => 'test'],
                    'taxId' => $taxId,
                    'visibilities' => [
                        [
                            'salesChannelId' => $salesChannelId,
                            'visibility' => ProductVisibilityDefinition::VISIBILITY_ALL
                        ],
                    ],
                    'categories' => [
                        ['id' => $categoryId],
                    ],
                ],
            ],
            Context::createDefaultContext()
        );
    }

    public function deactivate(DeactivateContext $deactivateContext): void
    {
        throw new \RuntimeException('It is not allowed to deactivate this Plugin');
    }

    public function uninstall(UninstallContext $uninstallContext): void
    {
        throw new \RuntimeException('It is not allowed to uninstall this Plugin');
    }

    public function build(ContainerBuilder $container): void
    {
        parent::build($container);

        $container->addCompilerPass(new MailTransportFactoryCompilerPass());
    }
}
