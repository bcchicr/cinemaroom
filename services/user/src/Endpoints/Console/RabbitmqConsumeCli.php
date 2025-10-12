<?php

declare(strict_types=1);

namespace App\Endpoints\Console;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;

#[AsCommand(
    name: 'rabbitmq:consume',
    description: 'start rabbitmq consumer',
)]
final readonly class RabbitmqConsumeCli
{
    public function __construct(
    ) {
    }

    public function __invoke(): int
    {
        echo 'Starting rabbitmq consumer...'.PHP_EOL;

        return Command::SUCCESS;
    }
}
