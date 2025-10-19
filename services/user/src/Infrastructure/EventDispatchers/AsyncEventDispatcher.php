<?php

declare(strict_types=1);

namespace App\Infrastructure\EventDispatchers;

use App\Domain\Events\DomainEvent;
use App\Domain\Events\EventDispatcher;

final readonly class AsyncEventDispatcher implements EventDispatcher
{
    #[\Override]
    public function dispatch(DomainEvent $event): void
    {
    }
}
