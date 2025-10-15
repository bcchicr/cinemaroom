<?php

declare(strict_types=1);

namespace App\Domain\ValueObjects;

use App\Domain\Abstracts\ValueObject;
use Ramsey\Uuid\UuidInterface;

abstract class Id extends ValueObject
{
    protected UuidInterface $value;

    public function __construct(
        UuidInterface $value,
    ) {
        $this->value = $value;
    }

    public function equals(ValueObject $other): bool
    {
        return $other instanceof static
            && $this->value->equals($other->value);
    }

    public function value(): UuidInterface
    {
        return $this->value;
    }

    public function hash(): string
    {
        return md5($this->toString());
    }

    public function toString(): string
    {
        return $this->value->toString();
    }

    public function __toString(): string
    {
        return $this->toString();
    }

    public function jsonSerialize(): array
    {
        return $this->toArray();
    }

    public function toArray(): array
    {
        return [
            'value' => $this->value->toString(),
        ];
    }
}
