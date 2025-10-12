<?php

declare(strict_types=1);

namespace App\Infrastructure\EventDispatchers;

use App\Domain\Events\DomainEvent;
use App\Domain\Events\EventDispatcher;

final class SyncEventDispatcher implements EventDispatcher
{
    private array $listeners = [];

    public function addListener(string $eventName, callable $listener): void
    {
        if (!isset($this->listeners[$eventName])) {
            $this->listeners[$eventName] = [];
        }
        $this->listeners[$eventName][] = $listener;
    }

    public function dispatch(DomainEvent $event): void
    {
        $eventName = $event->eventName();
        if (!isset($this->listeners[$eventName])) {
            return;
        }

        foreach ($this->listeners[$eventName] as $listener) {
            $listener($event);
        }
    }
}
