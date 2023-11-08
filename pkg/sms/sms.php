<?php

namespace app\models;


use yii\helpers\Json;

class Sms
{
  private static $login = "mixburger";
  private static $password = "CkcB6I8@z";
  private static $url = "http://91.204.239.44/broker-api/send";
  private static $port = "8083";
  private static $sender = "3700";

  public static function send($phone, $text)
  {
    $id = Dashboard::getGUID();
    $id = "mxb" . substr($id, 0, 10);
    $data = [];
    $data['messages'][] = [
      'recipient' => $phone,
      'message-id' => $id,
      'sms' => [
        'originator' => self::$sender,
        'content' => [
          'text' => $text
        ],
      ],
    ];
    $json = Json::encode($data);
    $curl = curl_init();
    $base = base64_encode(self::$login . ":" . self::$password);

    curl_setopt_array($curl, array(
//            CURLOPT_PORT => self::$port,
      CURLOPT_URL => self::$url,
      CURLOPT_RETURNTRANSFER => true,
      CURLOPT_ENCODING => "",
      CURLOPT_MAXREDIRS => 10,
      CURLOPT_TIMEOUT => 30,
      CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
      CURLOPT_CUSTOMREQUEST => "POST",
      CURLOPT_POSTFIELDS => $json,
      CURLOPT_HTTPHEADER => array(
        "Authorization: Basic {$base}",
        "Cache-Control: no-cache",
        "Content-Type: application/json",
      ),
    ));
    $response = curl_exec($curl);
    $err = curl_error($curl);
    curl_close($curl);
    if ($err) {
      return $err;
    }

    return $response;
  }

}