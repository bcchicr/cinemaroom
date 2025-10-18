<?php

declare(strict_types=1);

namespace App\Domain\Abstracts;

abstract class ValueObject implements \Stringable, \JsonSerializable
{
    abstract public function equals(ValueObject $other): bool;

    public function hash(): string
    {
        return hash('sha256', $this->toString());
    }

    abstract public function toString(): string;

    public function jsonSerialize(): array
    {
        return $this->toArray();
    }

    abstract public function toArray(): array;

    public function __toString(): string
    {
        return $this->toString();
    }
}
