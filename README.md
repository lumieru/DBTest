# DBTest
总体来说，postgresql效率比mongodb好。fasthttp和net/http没有发现明显的区别。postsql在预先prepare statement的情况下效率最好；然后是完全不prepare，自己构造sql语句发给数据库；最后是每次执行sql，先prepare，然后execute，这个效率最差。
