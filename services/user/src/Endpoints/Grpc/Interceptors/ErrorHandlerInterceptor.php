<?php

declare(strict_types=1);

namespace App\Endpoints\Grpc\Interceptors;

use App\Application\Exceptions\ApplicationException;
use App\Domain\Exceptions\DomainException;
use App\Endpoints\Exceptions\EndpointsException;
use App\Infrastructure\Attributes\GrpcStatusProvider\GrpcStatus;
use App\Infrastructure\Exceptions\InfrastructureException;
use Baldinof\RoadRunnerBundle\Grpc\InterceptorInterface;
use Baldinof\RoadRunnerBundle\RoadRunnerBridge\GrpcRequest;
use Baldinof\RoadRunnerBundle\RoadRunnerBridge\GrpcRequestInvokerInterface;
use Google\Protobuf\Any;
use GRPC\Services\Common\v1\CustomErrorDetails;
use Psr\Log\LoggerInterface;
use Spiral\RoadRunner\GRPC\Exception\GRPCException;

use const Grpc\STATUS_INTERNAL;

final readonly class ErrorHandlerInterceptor implements InterceptorInterface
{
    private const string STATUS_CODE_PREFIX = 'cinemaroom_users_';

    public function __construct(
        private LoggerInterface $logger,
    ) {}

    public function intercept(GrpcRequest $invocation, GrpcRequestInvokerInterface $next): \Iterator
    {
        $this->logger->info('gRPC method called', [
            'service' => $invocation->getService(),
            'method' => $invocation->getMethod(),
        ]);

        try {
            yield $next->invoke($invocation);

            $this->logger->info('gRPC response sent');
        } catch (\Throwable $exception) {
            $grpcErrorCode = $this->getGrpcErrorCode($exception);
            $customErrorCode = $this->getCustomErrorCode($exception);

            $this->logger->error('gRPC unexpected exception caught', [
                'message' => $exception->getMessage(),
                'grpcErrorCode' => $grpcErrorCode,
                'customErrorCode' => $customErrorCode,
                'trace' => $exception->getTraceAsString(),
            ]);

            $details = new CustomErrorDetails();
            $details->setCode($customErrorCode);
            $details->setMessage($exception->getMessage());

            throw new GRPCException($exception->getMessage(), $grpcErrorCode, [$details], $exception);
        }
    }

    public function getGrpcErrorCode(\Throwable $e): int
    {
        $reflection = new \ReflectionClass($e);
        $attributes = $reflection->getAttributes(GrpcStatus::class);
        if (!empty($attributes)) {
            $grpcStatus = $attributes[0]->newInstance();
            assert($grpcStatus instanceof GrpcStatus);

            return $grpcStatus->code;
        }

        return STATUS_INTERNAL;
    }

    private function getCustomErrorCode(\Throwable $e): string
    {
        if ($e instanceof DomainException) {
            return self::STATUS_CODE_PREFIX . 'domain_' . $e->getConventionalCode()->value;
        }

        if ($e instanceof ApplicationException) {
            return self::STATUS_CODE_PREFIX . 'application_' . $e->getConventionalCode()->value;
        }

        if ($e instanceof InfrastructureException) {
            return self::STATUS_CODE_PREFIX . 'infrastructure_' . $e->getConventionalCode()->value;
        }

        if ($e instanceof EndpointsException) {
            return self::STATUS_CODE_PREFIX . 'endpoints_' . $e->getConventionalCode()->value;
        }

        return self::STATUS_CODE_PREFIX . 'unknown';
    }
}
