<?php

declare(strict_types=1);

namespace App\Domain\Exceptions;

enum ErrorCode: string
{
    case FailedInvariant = 'failed_invariant';
    case InvalidArgument = 'invalid_argument';
}
