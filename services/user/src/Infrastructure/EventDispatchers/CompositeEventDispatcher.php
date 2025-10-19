<?php

declare(strict_types=1);

namespace App\Infrastructure\EventDispatchers;

use App\Domain\Events\DomainEvent;
use App\Domain\Events\EventDispatcher;

final readonly class CompositeEventDispatcher implements EventDispatcher
{
    public function __construct(
        private SyncEventDispatcher $syncEventDispatcher,
        private AsyncEventDispatcher $asyncEventDispatcher,
    ) {
    }

    #[\Override]
    public function dispatch(DomainEvent $event): void
    {
        $this->syncEventDispatcher->dispatch($event);
        $this->asyncEventDispatcher->dispatch($event);
    }
}
