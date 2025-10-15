<?php

declare(strict_types=1);

namespace App\Endpoints\Exceptions;

abstract class EndpointsException extends \RuntimeException
{
    abstract public function getConventionalCode(): ErrorCode;
}
