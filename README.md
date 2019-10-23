wxBot4g �ǻ���go��΢�Ż�����

����Ŀ����ά�����£���ӭstar��fork�� pull requests�� issue

## ��Դ
[wxBot](https://github.com/liuwons/wxBot)��һ���ǳ�����Ŀ�Դ΢�Ÿ��˺Žӿڣ�ʹ��Python���Կ�����[wxBot4g](https://github.com/liuwons/wxBot)��[wxBot](https://github.com/liuwons/wxBot)����go�İ汾��

## ��Ŀ��;
- �Ѹ���΢�ź���չΪ�����"������"���Զ��ظ�����Ů���ѣ�justСcase��
- ���������Լ�����Ŀ�У�ΪӦ���ṩ΢�Ż����˵�������
- ��¼�Լ��ͷ�Ⱥ������������ÿ���������������ÿ��������
- ���������豸������ʵ�֡�����������

## �ص�
- 1�������á�ʵ����Ϣ�ص�������ʵ����wcbot��wcbot.run()���ɡ�
- 2���ȶ������ߡ��ײ�ʵ���ˡ�������ά��΢�ŻỰ���ƣ����û�͸����
- 3���ṩrestful api��֧�ַ�����Ϣ��ָ��Ⱥ��ָ�����ѡ�
- 4���������ٶ���ࡣ

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
- [ ] �ṩrestful api��������Ϣ��ָ������/Ⱥ
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


## �ο�
- [�ھ�΢��Web��ͨ�ŵ�ȫ����](http://www.tanhao.me/talk/1466.html/)
- [Python��ҳ΢��API](https://github.com/liuwons/wxBot)
- [΢�Ÿ��˺Ż�����](https://github.com/newflydd/itchat4go)
