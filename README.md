wxBot4g �ǻ���go��΢�Ż�����

## ����
- gin��http��ܣ�
- cron����ʱ����
- etree������xml��
- viper�������ļ���ȡ��
- logrus����־��ܣ�
- go-qrcode����½��ά�����ɣ�

## Ŀǰ֧�ֵ���Ϣ����
### ������Ϣ
- [x] �ı�
- [x] ͼƬ
- [x] ����λ��
- [x] ������Ƭ
- [x] ����
- [x] С��Ƶ
- [ ] ����

### Ⱥ��Ϣ
- [x] �ı�
- [x] ͼƬ
- [x] ����λ��
- [x] ������Ƭ
- [x] ����
- [ ] ����

### TODO����
- [x] �ṩrestful api��������Ϣ��ָ������/Ⱥ
- [ ] �ļ�/ͼƬ�ϴ�������oss
- [ ] ����ָ��Ⱥ����
- [ ] �����¼���ķ�������з���

## ʹ������
24�д����ʵ��΢�Ż����˵ļ�����Ϣ����

    package main

    import (
        "wxBot4g/models"
        "wxBot4g/pkg/define"
        "wxBot4g/wcbot"
    
        "github.com/sirupsen/logrus"
    )
    
    func HandleMsg(msg models.RealRecvMsg) {
        logrus.Debug("MsgType: ", msg.MsgType, " ", " MsgTypeId: ", msg.MsgTypeId)
        logrus.Info(
            "��Ϣ����:", define.MsgIdString(msg.MsgType), " ",
            "��������:", define.MsgTypeIdString(msg.MsgTypeId), " ",
            "������:", msg.SendMsgUSer.Name, " ",
            "����:", msg.Content.Data)
    }
    
    func main() {
        bot := wcbot.New(HandleMsg)
        bot.Debug = true
        bot.Run()
    }


## ��Ϣ���ͺ���������

### MsgType����Ϣ���ͣ�

�������ͱ��|��������|˵��
--|--|--|
0|Init|��ʼ����Ϣ���ڲ�����
1|Self|�Լ����͵���Ϣ
2|FileHelper|�ļ���Ϣ
3|Group|Ⱥ��Ϣ
4|Contact|��ϵ����Ϣ
5|Public|���ں���Ϣ
6|Special|�����˺���Ϣ
51|��ȡwxid|��ȡwxid��Ϣ
99|Unknown|δ֪�˺���Ϣ


### MsgTypeId���������ͣ�

�������ͱ��|��������|˵��
--|--|--|
0|Text|�ı���Ϣ�ľ�������
1|Location|����λ��
3|Image|ͼƬ���ݵ�url��HTTP POST�����url���Եõ�jpg�ļ���ʽ������
4|Voice|�������ݵ�url��HTTP POST�����url���Եõ�mp3�ļ���ʽ������
5|Recommend|���� nickname (�ǳ�)�� alias (����)��province (ʡ��)��city (����)�� gender (�Ա�)�ֶ�
6|Animation|����url, HTTP POST�����url���Եõ�gif�ļ���ʽ������
7|Share|�ֵ䣬���� type (����)��title (����)��desc (����)��url (����)��from (Դ��վ)�ֶ�
8|Video|��Ƶ��δ֧��
9|VideoCall|��Ƶ�绰��δ֧��
10|Redraw|������Ϣ
11|Empty|���ݣ�δ֧��
99|Unknown|δ֧��

## ����api
### �����ı���Ϣ(����/Ⱥ)
```http://127.0.0.1:7788/v1/msg/text?to=����Ⱥ&word=���, ����һ��&appKey=khr1244o1oh```

### ����ͼƬ��Ϣ(����/Ⱥ)
��ο�```wxBot4g/wcbot/imageHandle_test.go```

v1.1
- �����ն˶�ά��ɨ���¼
- ����api�������ı���ͼƬ��Ϣ��ָ��Ⱥ
- ���ӵ�Ԫ����