<?php

namespace SBP\DemoServer\Controller;

use Docker\API\Exception\ContainerKillNotFoundException;
use Slim\Http\Request;
use Slim\Http\Response;

class Delete extends AbstractController
{
    public function delete(Request $request, Response $response)
    {
        $id = $request->getQueryParam('id');

        try {
            $this->getDocker()->delete($id);
        } catch (ContainerKillNotFoundException $e) {
            return $response
                ->withJson(['success' => true]);
        }

        return $response
                ->withJson(['success' => true]);
    }
}
