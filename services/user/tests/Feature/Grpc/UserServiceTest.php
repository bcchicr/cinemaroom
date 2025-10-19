<?php

declare(strict_types=1);

namespace App\Tests\Feature\Grpc;

use Google\Protobuf\Any;
use Google\Rpc\Status;
use Grpc\Channel;
use GRPC\Services\Common\v1\CustomErrorDetails;
use GRPC\Services\Common\v1\Uuid;
use GRPC\Services\Users\v1\Avatar;
use GRPC\Services\Users\v1\EditRequest;
use GRPC\Services\Users\v1\EditResponse;
use GRPC\Services\Users\v1\GetRequest;
use GRPC\Services\Users\v1\GetResponse;
use GRPC\Services\Users\v1\RegisterRequest;
use GRPC\Services\Users\v1\RegisterResponse;
use GRPC\Services\Users\v1\UserServiceInterface;

use const Grpc\STATUS_INVALID_ARGUMENT;
use const Grpc\STATUS_NOT_FOUND;
use const Grpc\STATUS_OK;

use Grpc\UnaryCall;
use PHPUnit\Framework\Attributes\DataProvider;
use PHPUnit\Framework\TestCase;

final class UserServiceTest extends TestCase
{
    private const ERROR_CODE_PREFIX = 'cinemaroom_users_';
    private Channel $channel;

    public static function invalid_username_provider(): array
    {
        return [
            [str_repeat('a', 256)],
            [''],
            [' '],
            ["\t"],
            ["\n"],
            ["\r"],
        ];
    }

    public function testRegister(): void
    {
        $username = 'GRPC_TEST_USER';

        $call = new UnaryCall(
            $this->channel,
            UserServiceInterface::NAME.'/Register',
            [RegisterResponse::class, 'mergeFromString'],
        );

        $request = new RegisterRequest();
        $request->setUsername($username);

        $call->start($request);
        [$response, $status] = $call->wait();

        $this->assertEquals(STATUS_OK, $status->code);
        $this->assertInstanceOf(RegisterResponse::class, $response);
        $this->assertEquals($response->getUser()->getUsername(), $username);
        $this->assertEquals('', $response->getUser()->getBio());
        $this->assertFalse($response->getUser()->hasAvatar());
    }

    /**
     * @throws \Exception
     */
    #[DataProvider('invalid_username_provider')]
    public function testRegisterWithInvalidUsername(string $username): void
    {
        $call = new UnaryCall(
            $this->channel,
            UserServiceInterface::NAME.'/Register',
            [RegisterResponse::class, 'mergeFromString'],
        );

        $request = new RegisterRequest();
        $request->setUsername($username);

        $call->start($request);
        [$response, $status] = $call->wait();

        $this->assertEquals(STATUS_INVALID_ARGUMENT, $status->code);
        $this->assertNull($response);
        $errorDetails = $this->unpackErrorDetails($status);

        $this->assertStringStartsWith(self::ERROR_CODE_PREFIX, $errorDetails->getCode());
    }

    /**
     * @throws \Exception
     */
    private function unpackErrorDetails(object $status): CustomErrorDetails
    {
        $statusDetailsBin = $status->metadata['grpc-status-details-bin'][0];

        $status = new Status();
        $status->mergeFromString($statusDetailsBin);

        $details = $status->getDetails();

        /**
         * @var Any $any
         */
        $any = $details[0];

        $errorDetails = $any->unpack();
        assert($errorDetails instanceof CustomErrorDetails);

        return $errorDetails;
    }

    public function testGetExistingUser(): void
    {
        $call = new UnaryCall(
            $this->channel,
            UserServiceInterface::NAME.'/Register',
            [RegisterResponse::class, 'mergeFromString'],
        );

        $request = new RegisterRequest();
        $request->setUsername('test_user');

        $call->start($request);
        [$response, $status] = $call->wait();
        $this->assertInstanceOf(RegisterResponse::class, $response);

        $id = $response->getUser()->getId();

        $call = new UnaryCall(
            $this->channel,
            UserServiceInterface::NAME.'/Get',
            [GetResponse::class, 'mergeFromString'],
        );
        $request = new GetRequest();
        $request->setId($id);
        $call->start($request);
        [$response, $status] = $call->wait();

        $this->assertEquals(STATUS_OK, $status->code);
        $this->assertInstanceOf(GetResponse::class, $response);
        $this->assertEquals($id->getValue(), $response->getUser()->getId()->getValue());
    }

    /**
     * @throws \Exception
     */
    public function testGetNonExistingUser(): void
    {
        $id = new Uuid();
        $id->setValue(\Ramsey\Uuid\Uuid::NIL);

        $call = new UnaryCall(
            $this->channel,
            UserServiceInterface::NAME.'/Get',
            [GetResponse::class, 'mergeFromString'],
        );
        $request = new GetRequest();
        $request->setId($id);
        $call->start($request);
        [$response, $status] = $call->wait();

        $this->assertEquals(STATUS_NOT_FOUND, $status->code);
        $this->assertNull($response);

        $errorDetails = $this->unpackErrorDetails($status);
        $this->assertStringStartsWith(self::ERROR_CODE_PREFIX, $errorDetails->getCode());
    }

    public function testEditExistingUser(): void
    {
        $call = new UnaryCall(
            $this->channel,
            UserServiceInterface::NAME.'/Register',
            [RegisterResponse::class, 'mergeFromString'],
        );

        $oldUsername = 'test_user';
        $request = new RegisterRequest();
        $request->setUsername($oldUsername);

        $call->start($request);
        [$response, $status] = $call->wait();
        $this->assertInstanceOf(RegisterResponse::class, $response);

        $id = $response->getUser()->getId();

        $call = new UnaryCall(
            $this->channel,
            UserServiceInterface::NAME.'/Edit',
            [EditResponse::class, 'mergeFromString'],
        );

        $newUsername = 'new_username';
        $newBio = 'new_bio';
        $newAvatarPath = '/path/to/avatar.jpg';
        $request = new EditRequest();
        $request->setId($id);
        $request->setUsername($newUsername);
        $request->setBio($newBio);
        $request->setAvatar((new Avatar())->setPath($newAvatarPath));

        $call->start($request);
        [$response, $status] = $call->wait();

        $this->assertEquals(STATUS_OK, $status->code);
        $this->assertInstanceOf(EditResponse::class, $response);
        $this->assertEquals($id->getValue(), $response->getUser()->getId()->getValue());
        $this->assertEquals($newUsername, $response->getUser()->getUsername());
        $this->assertEquals($newBio, $response->getUser()->getBio());
        $this->assertTrue($response->getUser()->hasAvatar());
        $this->assertEquals($newAvatarPath, $response->getUser()->getAvatar()->getPath());
    }

    /**
     * @throws \Exception
     */
    public function testEditNonExistingUser(): void
    {
        $call = new UnaryCall(
            $this->channel,
            UserServiceInterface::NAME.'/Edit',
            [EditResponse::class, 'mergeFromString'],
        );

        $newUsername = 'new_username';
        $newBio = 'new_bio';
        $newAvatarPath = '/path/to/avatar.jpg';
        $request = new EditRequest();
        $request->setId((new Uuid())->setValue(\Ramsey\Uuid\Uuid::NIL));
        $request->setUsername($newUsername);
        $request->setBio($newBio);
        $request->setAvatar((new Avatar())->setPath($newAvatarPath));

        $call->start($request);
        [$response, $status] = $call->wait();

        $this->assertEquals(STATUS_NOT_FOUND, $status->code);
        $this->assertNull($response);

        $errorDetails = $this->unpackErrorDetails($status);
        $this->assertStringStartsWith(self::ERROR_CODE_PREFIX, $errorDetails->getCode());
    }

    protected function setUp(): void
    {
        parent::setUp();

        $this->channel = new Channel(
            'localhost:9001',
            [
                'credentials' => \Grpc\ChannelCredentials::createInsecure(),
            ],
        );
    }
}
