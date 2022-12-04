INSERT INTO user (id,username,passwrd,email,created_at)
VALUES
    (1,"danilo","cangucu.danilo","danilo.cangucu@gritlab.ax",DateTime('now')),
    (2,"dang","lam.quoc.dang","dang.lam@gritlab.ax",DateTime('now')),
    (3,"iuliia","chipsanova.iuliia","iuliia.chipsanova@gritlab.ax",DateTime('now')),
    (4,"malin","oscarius.malin","malin.oscarius@gritlab.ax",DateTime('now')),
    (5,"tommy","mathisen.tommy","tommy.mathisen@gritlab.ax",DateTime('now'));

INSERT INTO category (id,category_name,created_at)
VALUES
    (1,"Cuisines",DateTime('now')),
    (2,"Places",DateTime('now')),
    (3,"Activities",DateTime('now'));

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
