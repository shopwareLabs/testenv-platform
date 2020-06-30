<?php

namespace SBP\DemoServer\Struct;

use Slim\Http\Request;

abstract class Struct implements \JsonSerializable
{
    public function fillFromRequest(Request $request): self
    {
        $body = $request->getParsedBody();

        if (is_array($body)) {
            foreach ($body as $k => $v) {
                if (property_exists($this, $k)) {
                    $this->$k = $v;
                }
            }
        }

        return $this;
    }

    public function jsonSerialize()
    {
        return get_object_vars($this);
    }
}