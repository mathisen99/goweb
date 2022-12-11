INSERT INTO user (id,privilege,username,passwrd,email,created_at)
VALUES
    (1,5,"danilo","$2a$14$DEdZfmpt910DZh41/H8ApeTecGFZ55v1p.DETh2/DXJ9D0ksqVkhu","danilo.cangucu@gritlab.ax",DateTime('now')),
    (2,5,"dang","$2a$14$GTXlBMDLMJq0UXCIq40FF.UOxqDP7aMVl/dh5SRdbSBl5Ss2aYWUe","dang.lam@gritlab.ax",DateTime('now')),
    (3,5,"iuliia","$2a$14$x/p9Ds7EUYxhX3jRaTOCKue9PByvhBatwsEIZx/3laO19vQcpiRdK","iuliia.chipsanova@gritlab.ax",DateTime('now')),
    (4,5,"malin","$2a$14$2DzMB74I10pRtzy3/5BHd.OXsj3S5iXYiq7ufALEMgm.NcCPzaZwK","malin.oscarius@gritlab.ax",DateTime('now')),
    (5,5,"tommy","$2a$14$isZC5HJtOGWx2EAInUqdF.GhQmZzoJRS7DQu767W5Z1aOU4.0ZP5G","tommy.mathisen@gritlab.ax",DateTime('now'));

INSERT INTO category (id,category_name,descript,created_at)
VALUES
    (1,"Cuisines","Recommendation regarding food in Mariehamn",DateTime('now')),
    (2,"Places","Places worth a visit in Mariehamn",DateTime('now')),
    (3,"Activities","Interesting events happening in Mariehamn",DateTime('now'));

INSERT INTO post (id,user_id,title,content,created_at,liked_no,disliked_no)
VALUES
    (1,2,"Asian Food","Thai Khun Mom serves very typical Asian food in Mariehamn",DateTime('now'),0,0),
    (2,3,"Swedish Class","Swedish class occurs every Tuesday and Thursday from 4pm",DateTime('now'),0,0),
    (3,4,"Best Sushi","Fina Fisken is the best sushi in Mariehamn",DateTime('now'),0,0),
    (4,5,"Poker Night","Poker Game Night occurs every Friday from 8pm",DateTime('now'),0,0),
    (5,1,"Real Embassy","Brazilian Real Embassy is now in Mariehamn",DateTime('now'),0,0);

INSERT INTO category_relation (id,category_id,post_id)
VALUES
    (1,1,1),
    (2,1,3),
    (3,2,1),
    (4,2,3),
    (5,2,5),
    (6,3,2),
    (7,3,4);
