import unittest

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from model import User, Post, Base

class TestOneToMany(unittest.TestCase):

    def setUp(self):
        self.engine = create_engine('sqlite://')
        DBSession = sessionmaker(bind=self.engine)
        self.session = DBSession()
        Base.metadata.create_all(self.engine)

    def tearDown(self):
        Base.metadata.drop_all(self.engine)

    def test_user_posts(self):
        u1 = User(name='user_one')
        u2 = User(name='user_two')

        self.session.add_all([u1, u2])
        self.session.commit()
        #No user post content, it's [].
        self.assertEqual(u1.posts, [])
        self.assertEqual(u2.posts, [])

        p1 = Post(content='This is belong to user_one', author=u1)
        p11 = Post(content='This is belong to user_one (2 posts)', author=u1)
        p2 = Post(content='This is belong to user_two', author=u2)
        p21 = Post(content='This is belong to user_two (2 posts)', author=u2)

        self.session.add_all([p1, p11, p2, p21])
        self.session.commit()
        #Now we have add post to user respectively.
        self.assertCountEqual(u1.posts, [p1, p11])
        self.assertCountEqual(u2.posts, [p2, p21])
        # Check post author.
        self.assertEqual(p1.author.name, 'user_one')
        self.assertEqual(p2.author.name, 'user_two')

if __name__ == '__main__':
    unittest.main()