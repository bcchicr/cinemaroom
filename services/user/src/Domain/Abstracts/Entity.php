<?php

declare(strict_types=1);

namespace App\Domain\Abstracts;

use App\Domain\ValueObjects\Id;

abstract class Entity
{
    public function equals(Entity $other): bool
    {
        return $this->id()->equals($other->id());
    }

    abstract public function id(): Id;
}
