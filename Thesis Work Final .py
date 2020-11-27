#!/usr/bin/env python
# coding: utf-8

# In[1]:


from pyspark import SparkContext, SQLContext, SparkConf
from pyspark.sql import SparkSession
from pyspark.sql import Row
from pyspark.sql import functions as F
from pyspark.sql.functions import from_utc_timestamp
import os
import requests


# In[3]:


# spark.stop()

conf = SparkConf().setAppName("PySpark App").setMaster("local").set("spark.jars", "/Users/altunbaratli/Desktop/Development/Robust-Data-Flow-Over-Blockchain").set('spark.driver.host','127.0.0.1')
sc = SparkContext(conf=conf)
sqlContext = SQLContext(sc)
spark = SparkSession     .builder     .appName("Python Spark SQL basic example").enableHiveSupport().getOrCreate()


# In[5]:


pload = {'username':'user1','channel':'channel1', 'smartcontract':'cc', 'args': {'sensorID':'1'}}
r = requests.post('http://169.51.200.103:30006/api/getHistory',json = pload)
print(r.text)


# In[6]:


json_rdd = sc.parallelize([r.text])
convertdf = spark.read.json(json_rdd)
convertdf.printSchema()


# In[13]:


newdf = convertdf.select("Value.sensorID", "Value.time", "Value.temp", "Value.stuckatfault", "Value.scfault", "Value.outlier")
flattened = newdf.dropna()
flattened_full = flattened.where("temp!=' '")
# df = flattened_full.withColumn('new_date',F.to_date(F.unix_timestamp('time', 'yyyy-MM-dd HH:mm:ss').cast('timestamp')))
df = flattened_full.withColumn('temp',F.col('temp').cast('float'))
plt = df.toPandas().plot(x="time", y="temp", figsize=(12,6), rot=50)


# In[ ]:




