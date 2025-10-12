<?php

declare(strict_types=1);

namespace App\Domain\Abstracts;

abstract class ValueObject implements \Stringable, \JsonSerializable
{
    abstract public function equals(ValueObject $other): bool;

    abstract public function toArray(): array;

    abstract public function toString(): string;

    abstract public function hash(): string;
}
