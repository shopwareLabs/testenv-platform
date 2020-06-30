<?php

namespace SBP\DemoServer\Services;

use Docker\API\Model\ContainersCreatePostBody;
use Docker\API\Model\ContainersCreatePostBodyNetworkingConfig;
use Docker\API\Model\ContainerSummaryItem;
use Docker\API\Model\HostConfig;
use Docker\API\Model\NetworkSettings;
use SBP\DemoServer\Struct\CreateShopRequest;

class Docker
{
    /**
     * @var \Docker\Docker
     */
    private $docker;

    public function __construct(\Docker\Docker $docker)
    {
        $this->docker = $docker;
    }

    public function createDemoShop(CreateShopRequest $request): array
    {
        $id = 'e' . random_int(1, 999);

        list($volumeFolder, $pluginName) = $this->preparePluginPath($request, $id);
        $domain = strtolower($pluginName .'-' .  $id. '.' . $_SERVER['BASE_HOST']);

        $body = new ContainersCreatePostBody();
        $body->setImage('shopware/testenv:' . $request->installVersion);
        $body->setEnv([
            'VIRTUAL_HOST=' . $domain,
            'PLUGIN_NAME=' . $pluginName
        ]);

        $hostConfig = new HostConfig();
        $hostConfig->setBinds([
            $volumeFolder . '/old:/var/www/shop/engine/Shopware/Plugins/Community/',
            $volumeFolder . '/new:/var/www/shop/custom/plugins/',
        ]);
        $body->setHostConfig($hostConfig);

        $networkSettings = new ContainersCreatePostBodyNetworkingConfig();
        $networkSettings->setEndpointsConfig(new \ArrayObject(['docker_default' => new NetworkSettings()]));
        $body->setNetworkingConfig($networkSettings);
        $body->setLabels(new \ArrayObject([
            'traefik.frontend.rule' => sprintf('Host: %s', $domain)
        ]));

        $this->docker->containerCreate($body, ['name' => $id]);
        $this->docker->containerStart($id);

        return [
            'id' => $id,
            'domain' => $domain . '/shop/public',
            'installedVersion' => $request->installVersion
        ];
    }

    public function list()
    {
        $response = $this->docker->containerList([
            'all' => true,
            'filters' => json_encode(['name' => ['e']])
        ]);

        $list = [];

        /** @var ContainerSummaryItem $item */
        foreach ($response as $item) {
            $list[] = [
                'id' => substr($item->getNames()[0], 1),
                'domain' => substr($item->getLabels()['traefik.frontend.rule'], 6),
                'installedVersion' => explode(':', $item->getImage())[1],
                'state' => $item->getState(),
                'status' => $item->getStatus()
            ];
        }

        return $list;
    }

    public function delete(string $id): void
    {
        $this->docker->containerKill($id);
        $this->docker->containerDelete($id);
    }

    private function preparePluginPath(CreateShopRequest $request, string $id): array
    {
        $volumeFolder = '/tmp/' . uniqid($id, true);
        mkdir($volumeFolder);
        mkdir($volumeFolder . '/new');
        mkdir($volumeFolder . '/old/');
        mkdir($volumeFolder . '/old/Frontend');
        mkdir($volumeFolder . '/old/Core');
        mkdir($volumeFolder . '/old/Backend');

        $zipUnique = '/tmp/' . uniqid('lel', true);
        file_put_contents($zipUnique, base64_decode($request->plugin));

        $zip = new \ZipArchive();
        $zip->open($zipUnique);

        $entry = $zip->statIndex(0);
        $rootDirectory = explode('/', $entry['name'])[0];

        $legacyPlugin = in_array($rootDirectory, ['Frontend', 'Backend', 'Core']);
        if ($legacyPlugin) {
            $zip->extractTo($volumeFolder . '/old/');
        } else {
            $zip->extractTo($volumeFolder . '/new/');
        }

        $zip->close();
        $segment = explode('/', $entry['name']);

        exec('chown -R 1000:1000 ' . $volumeFolder);

        return [
            $volumeFolder,
            $legacyPlugin ? $segment[1] : $segment[0]
        ];
    }
}
