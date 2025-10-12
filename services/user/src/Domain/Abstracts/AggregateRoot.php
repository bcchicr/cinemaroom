<?php

declare(strict_types=1);

namespace App\Domain\Abstracts;

use App\Domain\Events\DomainEvent;

abstract class AggregateRoot extends Entity
{
    private array $pendingEvents = [];

    public function pullPendingEvents(): array
    {
        $events = $this->pendingEvents;
        $this->pendingEvents = [];

        return $events;
    }

    protected function record(DomainEvent $event): void
    {
        $this->pendingEvents[] = $event;
    }
}
