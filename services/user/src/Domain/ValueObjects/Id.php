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

    #[\Override]
    public function equals(ValueObject $other): bool
    {
        return $other instanceof static
            && $this->value->equals($other->value);
    }

    public function value(): UuidInterface
    {
        return $this->value;
    }

    #[\Override]
    public function toArray(): array
    {
        return [
            'value' => $this->value->toString(),
        ];
    }

    #[\Override]
    public function toString(): string
    {
        return $this->value->toString();
    }
}
