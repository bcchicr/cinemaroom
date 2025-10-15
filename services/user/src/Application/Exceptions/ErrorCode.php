<?php

declare(strict_types=1);

namespace App\Application\Exceptions;

enum ErrorCode: string
{
    case NotFound = 'NOT_FOUND';
}
